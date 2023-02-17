//go:build windows
// +build windows

package forkexec

import (
	"unsafe"
	_ "unsafe"
)

////go:linkname _STARTUPINFOEXW syscall._STARTUPINFOEXW
//type _STARTUPINFOEXW struct

//go:linkname rtlGetNtVersionNumbers syscall.rtlGetNtVersionNumbers
func rtlGetNtVersionNumbers(majorVersion *uint32, minorVersion *uint32, buildNumber *uint32)

//go:linkname newProcThreadAttributeList syscall.newProcThreadAttributeList
func newProcThreadAttributeList(maxAttrCount uint32) (*_PROC_THREAD_ATTRIBUTE_LIST, error)

//go:linkname deleteProcThreadAttributeList syscall.deleteProcThreadAttributeList
func deleteProcThreadAttributeList(attrlist *_PROC_THREAD_ATTRIBUTE_LIST)

//go:linkname updateProcThreadAttribute syscall.updateProcThreadAttribute
func updateProcThreadAttribute(attrlist *_PROC_THREAD_ATTRIBUTE_LIST, flags uint32, attr uintptr, value unsafe.Pointer, size uintptr, prevvalue unsafe.Pointer, returnedsize *uintptr) (err error)
