// +build linux darwin

package executor

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)


type RLimit struct {
	Which int
	RLim  syscall.Rlimit
}

// 打开并获取文件的描述符
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

// 重映射文件描述符
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

// 硬件计时器
func setHardTimer(realTimeLimit int) error {
	var prealt ITimerVal
	prealt.ItInterval.TvSec = uint64(math.Floor(float64(realTimeLimit) / 1000.0))
	prealt.ItInterval.TvUsec = uint64(realTimeLimit % 1000 * 1000)
	prealt.ItValue.TvSec = prealt.ItInterval.TvSec
	prealt.ItValue.TvUsec = prealt.ItInterval.TvUsec
	_, _, err := syscall.RawSyscall(syscall.SYS_SETITIMER, ITimerReal, uintptr(unsafe.Pointer(&prealt)), 0)
	if err != 0 {
		return fmt.Errorf("system call setitimer() error: %s", err)
	}
	return nil
}

// fork调用
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

// 获取管道数据
func getPipe() ([]int, error) {
	var pipe = []int{0, 0}
	err := syscall.Pipe(pipe)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

// 设置资源限制 (setrlimit)
func setLimit(timeLimit, memoryLimit, realTimeLimit, fileSizeLimit int) error {

	// Set stack limit
	stack := uint64(memoryLimit * 1024)
	if runtime.GOOS == "darwin" { // WTF?! >= 65mb caused an operation not permitted!
		stack = uint64(65500 * 1024)
	}

	rlimits := []RLimit{
		// Set time limit: RLIMIT_CPU
		{
			Which: syscall.RLIMIT_CPU,
			RLim: getRLimitEntity(
				uint64(math.Ceil(float64(timeLimit)/1000.0)),
				uint64(math.Ceil(float64(timeLimit)/1000.0)),
			),
		},
		// Set memory limit: RLIMIT_DATA
		{
			Which: syscall.RLIMIT_DATA,
			RLim: getRLimitEntity(
				uint64(memoryLimit*1024),
				uint64(memoryLimit*1024),
			),
		},
		// Set memory limit: RLIMIT_AS
		{
			Which: syscall.RLIMIT_AS,
			RLim: getRLimitEntity(
				uint64(memoryLimit*1024*2),
				uint64(memoryLimit*1024*2+1024),
			),
		},
		// Set stack limit
		{
			Which: syscall.RLIMIT_STACK,
			RLim: getRLimitEntity(
				stack,
				stack,
			),
		},
		// Set file size limit: RLIMIT_FSIZE
		{
			Which: syscall.RLIMIT_FSIZE,
			RLim: getRLimitEntity(
				uint64(fileSizeLimit),
				uint64(fileSizeLimit),
			),
		},
	}

	for _, rlimit := range rlimits {
		err := syscall.Setrlimit(rlimit.Which, &rlimit.RLim)
		if err != nil {
			return fmt.Errorf("setrlimit(%d) error: %s", rlimit.Which, err)
		}
	}

	// Set time limit: setITimer
	if realTimeLimit > 0 {
		err := setHardTimer(realTimeLimit)
		if err != nil {
			return err
		}
	}

	return nil
}

func getRLimitEntity(cur, max uint64) syscall.Rlimit {
	return syscall.Rlimit{Cur: cur, Max: max}
}
