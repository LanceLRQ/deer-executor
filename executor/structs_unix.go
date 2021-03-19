// +build linux darwin

package executor

import (
	"github.com/LanceLRQ/deer-executor/v2/common/sandbox/process"
	"syscall"
)

// ProcessInfo 进程信息
type ProcessInfo struct {
	Pid     int                `json:"pid"`
	Process *process.Process   `json:"-"`
	Status  syscall.WaitStatus `json:"status"`
	Rusage  *syscall.Rusage    `json:"rusage"`
}
