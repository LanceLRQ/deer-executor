// +build linux darwin

package executor

import "syscall"

// 进程信息
type ProcessInfo struct {
	Pid    uintptr            `json:"pid"`
	Status syscall.WaitStatus `json:"status"`
	Rusage syscall.Rusage     `json:"rusage"`
}
