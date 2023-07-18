//go:build linux || darwin
// +build linux darwin

package executor

import (
	"github.com/LanceLRQ/deer-executor/v3/executor/sandbox/cmd"
	"syscall"
)

// ProcessInfo 进程信息
type ProcessInfo struct {
	Pid     int                `json:"pid"`
	Process *cmd.Process       `json:"-"`
	Status  syscall.WaitStatus `json:"status"`
	Rusage  *syscall.Rusage    `json:"rusage"`
}
