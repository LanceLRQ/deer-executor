//go:build linux
// +build linux

package forkexec

type (
	_C_short     int16
	_C_int       int32
	_C_long      int64
	_C_long_long int64
)

const (
	_AT_FDCWD            = -0x64
	_AT_REMOVEDIR        = 0x200
	_AT_SYMLINK_NOFOLLOW = 0x100
	_AT_EACCESS          = 0x200
)
