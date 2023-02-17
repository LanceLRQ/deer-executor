//go:build windows
// +build windows

package forkexec

import "syscall"

const (
	_PROC_THREAD_ATTRIBUTE_PARENT_PROCESS = 0x00020000
	_PROC_THREAD_ATTRIBUTE_HANDLE_LIST    = 0x00020002
	_EXTENDED_STARTUPINFO_PRESENT         = 0x00080000
)

type _PROC_THREAD_ATTRIBUTE_LIST struct {
	_ [1]byte
}

type _STARTUPINFOEXW struct {
	syscall.StartupInfo
	ProcThreadAttributeList *_PROC_THREAD_ATTRIBUTE_LIST
}
