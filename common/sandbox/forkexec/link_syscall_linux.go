//go:build linux
// +build linux

package forkexec

import (
	"syscall"
	_ "unsafe"
)

//go:linkname runtime_BeforeFork syscall.runtime_BeforeFork
func runtime_BeforeFork()

//go:linkname runtime_AfterFork syscall.runtime_AfterFork
func runtime_AfterFork()

//go:linkname runtime_AfterForkInChild syscall.runtime_AfterForkInChild
func runtime_AfterForkInChild()

//go:linkname formatIDMappings syscall.formatIDMappings
func formatIDMappings(idMap []syscall.SysProcIDMap) []byte

//go:linkname writeIDMappings syscall.writeIDMappings
func writeIDMappings(path string, idMap []syscall.SysProcIDMap) error

//go:linkname writeSetgroups syscall.writeSetgroups
func writeSetgroups(pid int, enable bool) error

//go:linkname rawSyscallNoError syscall.rawSyscallNoError
func rawSyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr)

//go:linkname rawVforkSyscall syscall.rawVforkSyscall
func rawVforkSyscall(trap, a1 uintptr) (r1 uintptr, err syscall.Errno)
