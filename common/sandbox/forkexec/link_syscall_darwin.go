//go:build darwin
// +build darwin

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

//go:linkname rawSyscall syscall.rawSyscall
func rawSyscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:linkname libc_fork_trampoline syscall.libc_fork_trampoline
func libc_fork_trampoline()

//go:linkname libc_setsid_trampoline syscall.libc_setsid_trampoline
func libc_setsid_trampoline()

//go:linkname libc_setpgid_trampoline syscall.libc_setpgid_trampoline
func libc_setpgid_trampoline()

//go:linkname libc_getpid_trampoline syscall.libc_getpid_trampoline
func libc_getpid_trampoline()

//go:linkname libc_ioctl_trampoline syscall.libc_ioctl_trampoline
func libc_ioctl_trampoline()

//go:linkname libc_chroot_trampoline syscall.libc_chroot_trampoline
func libc_chroot_trampoline()

//go:linkname libc_setgroups_trampoline syscall.libc_setgroups_trampoline
func libc_setgroups_trampoline()

//go:linkname libc_setgid_trampoline syscall.libc_setgid_trampoline
func libc_setgid_trampoline()

//go:linkname libc_setuid_trampoline syscall.libc_setuid_trampoline
func libc_setuid_trampoline()

//go:linkname libc_chdir_trampoline syscall.libc_chdir_trampoline
func libc_chdir_trampoline()

//go:linkname libc_dup2_trampoline syscall.libc_dup2_trampoline
func libc_dup2_trampoline()

//go:linkname libc_fcntl_trampoline syscall.libc_fcntl_trampoline
func libc_fcntl_trampoline()

//go:linkname libc_close_trampoline syscall.libc_close_trampoline
func libc_close_trampoline()

//go:linkname libc_execve_trampoline syscall.libc_execve_trampoline
func libc_execve_trampoline()

//go:linkname libc_write_trampoline syscall.libc_write_trampoline
func libc_write_trampoline()

//go:linkname libc_exit_trampoline syscall.libc_exit_trampoline
func libc_exit_trampoline()

//go:linkname libc_read_trampoline syscall.libc_read_trampoline
func libc_read_trampoline()

//go:linkname libc_setrlimit_trampoline syscall.libc_setrlimit_trampoline
func libc_setrlimit_trampoline()

//go:linkname fcntl syscall.fcntl
func fcntl(fd int, cmd int, arg int) (val int, err error)

//go:linkname ptrace1 syscall.ptrace1
func ptrace1(request int, pid int, addr uintptr, data uintptr) (err error)

//go:linkname readlen syscall.readlen
func readlen(fd int, buf *byte, nbuf int) (n int, err error)

//go:linkname fork syscall.fork
func fork() (pid int, err error)
