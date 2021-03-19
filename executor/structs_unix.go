// +build linux darwin

package executor

import (
	"github.com/LanceLRQ/deer-common/sandbox/process"
	"syscall"
)

// 进程信息
type ProcessInfo struct {
	Pid     int                `json:"pid"`
	Process *process.Process   `json:"-"`
	Status  syscall.WaitStatus `json:"status"`
	Rusage  *syscall.Rusage    `json:"rusage"`
}
