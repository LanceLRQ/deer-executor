// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fork, exec, wait, etc.

package forkexec

import (
	"fmt"
	"golang.org/x/sys/windows"
	"log"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

var ForkLock sync.RWMutex

// EscapeArg rewrites command line argument s as prescribed
// in https://msdn.microsoft.com/en-us/library/ms880421.
// This function returns "" (2 double quotes) if s is empty.
// Alternatively, these transformations are done:
//   - every back slash (\) is doubled, but only if immediately
//     followed by double quote (");
//   - every double quote (") is escaped by back slash (\);
//   - finally, s is wrapped with double quotes (arg -> "arg"),
//     but only if there is space or tab inside s.
func EscapeArg(s string) string {
	if len(s) == 0 {
		return `""`
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\', ' ', '\t':
			// Some escaping required.
			b := make([]byte, 0, len(s)+2)
			b = appendEscapeArg(b, s)
			return string(b)
		}
	}
	return s
}

// appendEscapeArg escapes the string s, as per escapeArg,
// appends the result to b, and returns the updated slice.
func appendEscapeArg(b []byte, s string) []byte {
	if len(s) == 0 {
		return append(b, `""`...)
	}

	needsBackslash := false
	hasSpace := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\':
			needsBackslash = true
		case ' ', '\t':
			hasSpace = true
		}
	}

	if !needsBackslash && !hasSpace {
		// No special handling required; normal case.
		return append(b, s...)
	}
	if !needsBackslash {
		// hasSpace is true, so we need to quote the string.
		b = append(b, '"')
		b = append(b, s...)
		return append(b, '"')
	}

	if hasSpace {
		b = append(b, '"')
	}
	slashes := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		default:
			slashes = 0
		case '\\':
			slashes++
		case '"':
			for ; slashes > 0; slashes-- {
				b = append(b, '\\')
			}
			b = append(b, '\\')
		}
		b = append(b, c)
	}
	if hasSpace {
		for ; slashes > 0; slashes-- {
			b = append(b, '\\')
		}
		b = append(b, '"')
	}

	return b
}

// makeCmdLine builds a command line out of args by escaping "special"
// characters and joining the arguments with spaces.
func makeCmdLine(args []string) string {
	var b []byte
	for _, v := range args {
		if len(b) > 0 {
			b = append(b, ' ')
		}
		b = appendEscapeArg(b, v)
	}
	return string(b)
}

//go:linkname createEnvBlock syscall.createEnvBlock
func createEnvBlock(envv []string) (*uint16, error)

func CloseOnExec(fd syscall.Handle) {
	syscall.SetHandleInformation(syscall.Handle(fd), syscall.HANDLE_FLAG_INHERIT, 0)
}

func SetNonblock(fd syscall.Handle, nonblocking bool) (err error) {
	return nil
}

// FullPath retrieves the full path of the specified file.
func FullPath(name string) (path string, err error) {
	p, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return "", err
	}
	n := uint32(100)
	for {
		buf := make([]uint16, n)
		n, err = syscall.GetFullPathName(p, uint32(len(buf)), &buf[0], nil)
		if err != nil {
			return "", err
		}
		if n <= uint32(len(buf)) {
			return syscall.UTF16ToString(buf[:n]), nil
		}
	}
}

func isSlash(c uint8) bool {
	return c == '\\' || c == '/'
}

func normalizeDir(dir string) (name string, err error) {
	ndir, err := FullPath(dir)
	if err != nil {
		return "", err
	}
	if len(ndir) > 2 && isSlash(ndir[0]) && isSlash(ndir[1]) {
		// dir cannot have \\server\share\path form
		return "", syscall.EINVAL
	}
	return ndir, nil
}

func volToUpper(ch int) int {
	if 'a' <= ch && ch <= 'z' {
		ch += 'A' - 'a'
	}
	return ch
}

func joinExeDirAndFName(dir, p string) (name string, err error) {
	if len(p) == 0 {
		return "", syscall.EINVAL
	}
	if len(p) > 2 && isSlash(p[0]) && isSlash(p[1]) {
		// \\server\share\path form
		return p, nil
	}
	if len(p) > 1 && p[1] == ':' {
		// has drive letter
		if len(p) == 2 {
			return "", syscall.EINVAL
		}
		if isSlash(p[2]) {
			return p, nil
		} else {
			d, err := normalizeDir(dir)
			if err != nil {
				return "", err
			}
			if volToUpper(int(p[0])) == volToUpper(int(d[0])) {
				return FullPath(d + "\\" + p[2:])
			} else {
				return FullPath(p)
			}
		}
	} else {
		// no drive letter
		d, err := normalizeDir(dir)
		if err != nil {
			return "", err
		}
		if isSlash(p[0]) {
			return FullPath(d[:2] + p)
		} else {
			return FullPath(d + "\\" + p)
		}
	}
}

type ProcAttr struct {
	Dir   string
	Env   []string
	Files []uintptr
	Sys   *SysProcAttr
}

type SysProcAttr struct {
	HideWindow                 bool
	CmdLine                    string // used if non-empty, else the windows command line is built by escaping the arguments passed to StartProcess
	CreationFlags              uint32
	Token                      syscall.Token               // if set, runs new process in the security context represented by the token
	ProcessAttributes          *syscall.SecurityAttributes // if set, applies these security attributes as the descriptor for the new process
	ThreadAttributes           *syscall.SecurityAttributes // if set, applies these security attributes as the descriptor for the main thread of the new process
	NoInheritHandles           bool                        // if set, each inheritable handle in the calling process is not inherited by the new process
	AdditionalInheritedHandles []syscall.Handle            // a list of additional handles, already marked as inheritable, that will be inherited by the new process
	ParentProcess              syscall.Handle              // if non-zero, the new process regards the process given by this handle as its parent process, and AdditionalInheritedHandles, if set, should exist in this parent process
}

var zeroProcAttr ProcAttr
var zeroSysProcAttr SysProcAttr

func StartProcess(argv0 string, argv []string, attr *ProcAttr) (pid int, handle uintptr, err error) {
	if len(argv0) == 0 {
		return 0, 0, syscall.EWINDOWS
	}
	if attr == nil {
		attr = &zeroProcAttr
	}
	sys := attr.Sys
	if sys == nil {
		sys = &zeroSysProcAttr
	}

	if len(attr.Files) > 3 {
		return 0, 0, syscall.EWINDOWS
	}
	if len(attr.Files) < 3 {
		return 0, 0, syscall.EINVAL
	}

	if len(attr.Dir) != 0 {
		// StartProcess assumes that argv0 is relative to attr.Dir,
		// because it implies Chdir(attr.Dir) before executing argv0.
		// Windows CreateProcess assumes the opposite: it looks for
		// argv0 relative to the current directory, and, only once the new
		// process is started, it does Chdir(attr.Dir). We are adjusting
		// for that difference here by making argv0 absolute.
		var err error
		argv0, err = joinExeDirAndFName(attr.Dir, argv0)
		if err != nil {
			return 0, 0, err
		}
	}
	argv0p, err := syscall.UTF16PtrFromString(argv0)
	if err != nil {
		return 0, 0, err
	}

	var cmdline string
	// Windows CreateProcess takes the command line as a single string:
	// use attr.CmdLine if set, else build the command line by escaping
	// and joining each argument with spaces
	if sys.CmdLine != "" {
		cmdline = sys.CmdLine
	} else {
		cmdline = makeCmdLine(argv)
	}

	var argvp *uint16
	if len(cmdline) != 0 {
		argvp, err = syscall.UTF16PtrFromString(cmdline)
		if err != nil {
			return 0, 0, err
		}
	}

	var dirp *uint16
	if len(attr.Dir) != 0 {
		dirp, err = syscall.UTF16PtrFromString(attr.Dir)
		if err != nil {
			return 0, 0, err
		}
	}

	var maj, min, build uint32
	rtlGetNtVersionNumbers(&maj, &min, &build)
	isWin7 := maj < 6 || (maj == 6 && min <= 1)
	// NT kernel handles are divisible by 4, with the bottom 3 bits left as
	// a tag. The fully set tag correlates with the types of handles we're
	// concerned about here.  Except, the kernel will interpret some
	// special handle values, like -1, -2, and so forth, so kernelbase.dll
	// checks to see that those bottom three bits are checked, but that top
	// bit is not checked.
	isLegacyWin7ConsoleHandle := func(handle syscall.Handle) bool { return isWin7 && handle&0x10000003 == 3 }

	p, _ := syscall.GetCurrentProcess()
	parentProcess := p
	if sys.ParentProcess != 0 {
		parentProcess = sys.ParentProcess
	}
	fd := make([]syscall.Handle, len(attr.Files))
	for i := range attr.Files {
		if attr.Files[i] > 0 {
			destinationProcessHandle := parentProcess

			// On Windows 7, console handles aren't real handles, and can only be duplicated
			// into the current process, not a parent one, which amounts to the same thing.
			if parentProcess != p && isLegacyWin7ConsoleHandle(syscall.Handle(attr.Files[i])) {
				destinationProcessHandle = p
			}

			err := syscall.DuplicateHandle(p, syscall.Handle(attr.Files[i]), destinationProcessHandle, &fd[i], 0, true, syscall.DUPLICATE_SAME_ACCESS)
			if err != nil {
				return 0, 0, err
			}
			defer syscall.DuplicateHandle(parentProcess, fd[i], 0, nil, 0, false, syscall.DUPLICATE_CLOSE_SOURCE)
		}
	}
	si := new(_STARTUPINFOEXW)
	si.ProcThreadAttributeList, err = newProcThreadAttributeList(2)
	if err != nil {
		return 0, 0, err
	}
	defer deleteProcThreadAttributeList(si.ProcThreadAttributeList)
	si.Cb = uint32(unsafe.Sizeof(*si))
	si.Flags = syscall.STARTF_USESTDHANDLES
	if sys.HideWindow {
		si.Flags |= syscall.STARTF_USESHOWWINDOW
		si.ShowWindow = syscall.SW_HIDE
	}
	if sys.ParentProcess != 0 {
		err = updateProcThreadAttribute(si.ProcThreadAttributeList, 0, _PROC_THREAD_ATTRIBUTE_PARENT_PROCESS, unsafe.Pointer(&sys.ParentProcess), unsafe.Sizeof(sys.ParentProcess), nil, nil)
		if err != nil {
			return 0, 0, err
		}
	}
	si.StdInput = fd[0]
	si.StdOutput = fd[1]
	si.StdErr = fd[2]

	fd = append(fd, sys.AdditionalInheritedHandles...)

	// On Windows 7, console handles aren't real handles, so don't pass them
	// through to PROC_THREAD_ATTRIBUTE_HANDLE_LIST.
	for i := range fd {
		if isLegacyWin7ConsoleHandle(fd[i]) {
			fd[i] = 0
		}
	}

	// The presence of a NULL handle in the list is enough to cause PROC_THREAD_ATTRIBUTE_HANDLE_LIST
	// to treat the entire list as empty, so remove NULL handles.
	j := 0
	for i := range fd {
		if fd[i] != 0 {
			fd[j] = fd[i]
			j++
		}
	}
	fd = fd[:j]

	willInheritHandles := len(fd) > 0 && !sys.NoInheritHandles

	// Do not accidentally inherit more than these handles.
	if willInheritHandles {
		err = updateProcThreadAttribute(si.ProcThreadAttributeList, 0, _PROC_THREAD_ATTRIBUTE_HANDLE_LIST, unsafe.Pointer(&fd[0]), uintptr(len(fd))*unsafe.Sizeof(fd[0]), nil, nil)
		if err != nil {
			return 0, 0, err
		}
	}

	envBlock, err := createEnvBlock(attr.Env)
	if err != nil {
		return 0, 0, err
	}

	pi := new(syscall.ProcessInformation)
	flags := sys.CreationFlags | syscall.CREATE_UNICODE_ENVIRONMENT | _EXTENDED_STARTUPINFO_PRESENT
	if sys.Token != 0 {
		err = syscall.CreateProcessAsUser(sys.Token, argv0p, argvp, sys.ProcessAttributes, sys.ThreadAttributes, willInheritHandles, flags, envBlock, dirp, &si.StartupInfo, pi)
	} else {
		err = syscall.CreateProcess(argv0p, argvp, sys.ProcessAttributes, sys.ThreadAttributes, willInheritHandles, flags, envBlock, dirp, &si.StartupInfo, pi)
	}
	if err != nil {
		return 0, 0, err
	}

	initJobObj(pi)
	fmt.Println("AA")

	defer syscall.CloseHandle(syscall.Handle(pi.Thread))
	runtime.KeepAlive(fd)
	runtime.KeepAlive(sys)

	return int(pi.ProcessId), uintptr(pi.Process), nil
}

func Exec(argv0 string, argv []string, envv []string) (err error) {
	return syscall.EWINDOWS
}

func initJobObj(pi *syscall.ProcessInformation) {
	jobObj, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	limit := windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
		LimitFlags:              windows.JOB_OBJECT_LIMIT_WORKINGSET | windows.JOB_OBJECT_LIMIT_PROCESS_TIME,
		MinimumWorkingSetSize:   1024 * 1024,
		MaximumWorkingSetSize:   1024 * 1024 * 10, // 10M
		PerProcessUserTimeLimit: 1000 * 1000 * 1,  //单位100纳秒,当前是1秒
	}
	_, err = windows.SetInformationJobObject(jobObj, windows.JobObjectBasicLimitInformation, uintptr(unsafe.Pointer(&limit)), uint32(unsafe.Sizeof(limit)))
	if err != nil {
		log.Fatalf("err: %v", err)
		return
	}

	hProcess := windows.Handle(pi.Process)
	err = windows.AssignProcessToJobObject(jobObj, hProcess)
	if err != nil {
		log.Fatalln(err)
		return
	}
	windows.ResumeThread(hProcess)
}
