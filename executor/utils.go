/* Deer executor
 * (C) 2019-Now LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package executor

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"syscall"
	"unsafe"
)

// 定义ITimer的常量，命名规则遵循Linux的原始设定
const (
	ITimerReal 		= 0
	ITimerVirtual 	= 1
	ITimerVProf 	= 2
)

type ITimerVal struct  {
	ItInterval TimeVal
	ItValue TimeVal
}

type TimeVal struct {
	TvSec uint64
	TvUsec uint64
}


// 打开文件并获取描述符 (open)
func openFile(filePath string, flag int, perm os.FileMode) (*os.File, error) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file (%s) not exists", filePath)
		} else {
			return nil, fmt.Errorf("open file (%s) error: %s", filePath, err.Error())
		}
	} else {
		if fp, err := os.OpenFile(filePath, flag, perm); err != nil {
			return nil, fmt.Errorf("open file (%s) error: %s", filePath, err.Error())
		} else {
			return fp, nil
		}
	}
}

// 设置定时器 (setitimer)
func setITimer(prealt ITimerVal) (err error) {
	_, _, errMsg := syscall.RawSyscall(syscall.SYS_SETITIMER, ITimerReal, uintptr(unsafe.Pointer(&prealt)), 0)
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
func setLimit(timeLimit, memoryLimit , realTimeLimit int) (err error) {
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

	if realTimeLimit > 0 {
		// Set time limit: setITimer
		prealt.ItInterval.TvSec = uint64(math.Floor(float64(realTimeLimit) / 1000.0))
		prealt.ItInterval.TvUsec = uint64(realTimeLimit % 1000 * 1000)
		prealt.ItValue.TvSec = prealt.ItInterval.TvSec
		prealt.ItValue.TvUsec = prealt.ItInterval.TvUsec
		errMsg = setITimer(prealt)
		if errMsg != nil {
			return errMsg
		}
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

// 文件读写
func readFile(filePath string, name string, tryOnFailed int) ([]byte, string, error) {
	errCnt, errText := 0, ""
	var err error
	for errCnt < tryOnFailed {
		fp, err := openFile(filePath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
		if err != nil {
			errText = err.Error()
			errCnt++
			continue
		}
		data, err := ioutil.ReadAll(fp)
		if err != nil {
			_ = fp.Close()
			errText = fmt.Sprintf("read %s file i/o error: %s", name, err.Error())
			errCnt++
			continue
		}
		_ = fp.Close()
		return data, errText, nil
	}
	return nil, errText, err
}