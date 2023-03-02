// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || linux
// +build darwin linux

// Fork, exec, wait, etc.

package forkexec

import (
	errorspkg "errors"
	"syscall"
	"unsafe"
)

// ProcAttr holds attributes that will be applied to a new process started
// by StartProcess.
type ProcAttr struct {
	Dir   string    // Current working directory.
	Env   []string  // Environment.
	Files []uintptr // File descriptors.
	Sys   *SysProcAttr
}

var zeroProcAttr ProcAttr
var zeroSysProcAttr SysProcAttr

func forkExec(argv0 string, argv []string, attr *ProcAttr) (pid int, err error) {
	var p [2]int
	var n int
	var err1 syscall.Errno
	var wstatus syscall.WaitStatus

	if attr == nil {
		attr = &zeroProcAttr
	}
	sys := attr.Sys
	if sys == nil {
		sys = &zeroSysProcAttr
	}

	// Convert args to C form.
	argv0p, err := syscall.BytePtrFromString(argv0)
	if err != nil {
		return 0, err
	}
	argvp, err := syscall.SlicePtrFromStrings(argv)
	if err != nil {
		return 0, err
	}
	envvp, err := syscall.SlicePtrFromStrings(attr.Env)
	if err != nil {
		return 0, err
	}

	var chroot *byte
	if sys.Chroot != "" {
		chroot, err = syscall.BytePtrFromString(sys.Chroot)
		if err != nil {
			return 0, err
		}
	}
	var dir *byte
	if attr.Dir != "" {
		dir, err = syscall.BytePtrFromString(attr.Dir)
		if err != nil {
			return 0, err
		}
	}

	// Both Setctty and Foreground use the Ctty field,
	// but they give it slightly different meanings.
	if sys.Setctty && sys.Foreground {
		return 0, errorspkg.New("both Setctty and Foreground set in SysProcAttr")
	}
	if sys.Setctty && sys.Ctty >= len(attr.Files) {
		return 0, errorspkg.New("Setctty set but Ctty not valid in child")
	}

	// Acquire the fork lock so that no other threads
	// create new fds that are not yet close-on-exec
	// before we fork.
	syscall.ForkLock.Lock()

	// Allocate child status pipe close on exec.
	if err = forkExecPipe(p[:]); err != nil {
		syscall.ForkLock.Unlock()
		return 0, err
	}

	// Kick off child.
	pid, err1 = forkAndExecInChild(argv0p, argvp, envvp, chroot, dir, attr, sys, p[1])
	if err1 != 0 {
		syscall.Close(p[0])
		syscall.Close(p[1])
		syscall.ForkLock.Unlock()
		return 0, err1
	}
	syscall.ForkLock.Unlock()

	// Read child error status from pipe.
	syscall.Close(p[1])
	for {
		n, err = readlen(p[0], (*byte)(unsafe.Pointer(&err1)), int(unsafe.Sizeof(err1)))
		if err != syscall.EINTR {
			break
		}
	}
	syscall.Close(p[0])
	if err != nil || n != 0 {
		if n == int(unsafe.Sizeof(err1)) {
			err = syscall.Errno(err1)
		}
		if err == nil {
			err = syscall.EPIPE
		}
		// Child failed; wait for it to exit, to make sure
		// the zombies don't accumulate.
		_, err1 := syscall.Wait4(pid, &wstatus, 0, nil)
		for err1 == syscall.EINTR {
			_, err1 = syscall.Wait4(pid, &wstatus, 0, nil)
		}
		return 0, err
	}

	// Read got EOF, so pipe closed on exec, so exec succeeded.
	return pid, nil
}

// ForkExec Combination of fork and exec, careful to be thread safe.
func ForkExec(argv0 string, argv []string, attr *ProcAttr) (pid int, err error) {
	return forkExec(argv0, argv, attr)
}

// StartProcess wraps ForkExec for package os.
func StartProcess(argv0 string, argv []string, attr *ProcAttr) (pid int, handle uintptr, err error) {
	pid, err = forkExec(argv0, argv, attr)
	return pid, 0, err
}
