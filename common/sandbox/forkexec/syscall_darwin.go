// +build darwin

package forkexec

import "syscall"

// GetPipe 获取管道数据
func GetPipe() ([]uintptr, error) {
	var pipe = []int{0, 0}
	err := syscall.Pipe(pipe)
	if err != nil {
		return nil, err
	}
	return []uintptr{uintptr(pipe[0]), uintptr(pipe[1])}, nil
}
