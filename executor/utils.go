/* Deer executor
 * (C) 2019-Now LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package executor

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

// 定义ITimer的常量，命名规则遵循Linux的原始设定
const (
	ITIMER_REAL 	= 0
	ITIMER_VIRTUAL 	= 1
	ITIMER_PROF 	= 2
)

type ITimerVal struct  {
	ItInterval TimeVal
	ItValue TimeVal
}

type TimeVal struct {
	TvSec uint64
	TvUsec uint64
}

// 创建一个新的进程 (fork)
func forkProc() (pid uintptr, err error) {
	r1, r2, errMsg := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	darwin := runtime.GOOS == "darwin"
	if errMsg != 0 {
		return 0, fmt.Errorf("system call: fork(); error: %s", errMsg)
	}
	if darwin {
		if r2 == 1 {
			pid = 0
		} else {
			pid = r1
		}
	} else {
		if r1 == 0 && r2 == 0 {
			pid = 0
		} else {
			pid = r1
		}
	}
	return pid, nil
}

// 打开文件并获取描述符 (open)
func getFileDescriptor(path string, flag int, perm uint32) (fd int, err error) {
	var filed = 0
	_, errMsg := os.Stat(path)
	if errMsg != nil {
		if os.IsNotExist(err) {
			return 0, errMsg
		}
	}
	filed, errMsg = syscall.Open(path, flag, perm)
	return filed, nil
}

// 重定向文件描述符 (dup2)
func redirectFileDescriptor(to int, path string, flag int, perm uint32) (fd int, err error) {
	fd, errMsg := getFileDescriptor(path, flag, perm)
	if errMsg == nil {
		errMsg = syscall.Dup2(fd, to)
		if errMsg != nil {
			syscall.Close(fd)
			return -1, errMsg
		}
		return fd, nil
	} else {
		return -1, errMsg
	}
}

// 设置定时器 (setitimer)
func setITimer(prealt ITimerVal) (err error) {
	_, _, errMsg := syscall.RawSyscall(syscall.SYS_SETITIMER, ITIMER_REAL, uintptr(unsafe.Pointer(&prealt)), 0)
	if errMsg != 0 {
		return fmt.Errorf("system call: setitimer(); error: %s", errMsg)
	}
	return nil
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// 设置资源限制 (setrlimit)
func setLimit(timeLimit int, memoryLimit int) (err error) {
	var rlimit syscall.Rlimit
	var prealt ITimerVal
	var errMsg error

	// Set time limit: RLIMIT_CPU
	rlimit.Cur = uint64(math.Ceil(float64(timeLimit) / 1000.0))
	rlimit.Max = rlimit.Cur + 1

	errMsg = syscall.Setrlimit(syscall.RLIMIT_CPU, &rlimit)
	if errMsg != nil {
		return errMsg
	}

	// Set time limit: setITimer
	prealt.ItInterval.TvSec = uint64(math.Floor(float64(timeLimit) / 1000.0))
	prealt.ItInterval.TvUsec = uint64(timeLimit % 1000 * 1000)
	prealt.ItValue.TvSec = prealt.ItInterval.TvSec
	prealt.ItValue.TvUsec = prealt.ItInterval.TvUsec
	errMsg = setITimer(prealt)
	if errMsg != nil {
		return errMsg
	}

	// Set memory limit: RLIMIT_DATA
	rlimit.Cur = uint64(memoryLimit * 1024)
	rlimit.Max = rlimit.Cur
	errMsg = syscall.Setrlimit(syscall.RLIMIT_DATA, &rlimit)
	if errMsg != nil {
		return errMsg
	}

	// Set memory limit: RLIMIT_AS
	rlimit.Cur = uint64(memoryLimit * 1024) * 2
	rlimit.Max = rlimit.Cur + 1024
	errMsg = syscall.Setrlimit(syscall.RLIMIT_AS, &rlimit)
	if errMsg != nil {
		return errMsg
	}

	// Set stack limit: RLIMIT_STACK
	rlimit.Cur = uint64(memoryLimit * 1024)
	rlimit.Max = rlimit.Cur
	errMsg = syscall.Setrlimit(syscall.RLIMIT_STACK, &rlimit)
	if errMsg != nil {
		return errMsg
	}

	// Set file size limit: RLIMIT_FSIZE
	rlimit.Cur = uint64(JudgeFileSizeLimit)
	rlimit.Max = rlimit.Cur
	errMsg = syscall.Setrlimit(syscall.RLIMIT_FSIZE, &rlimit)


	return nil
}
