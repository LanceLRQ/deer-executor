// +build darwin linux

package forkexec

import (
	"math"
	"runtime"
	"syscall"
)

const DarwinSafeStackSize = 65500

// 定义ITimer的常量
const (
	ITIMER_REAL    = 0
	ITIMER_VIRTUAL = 1
	ITIMER_PROF    = 2
)

type RLimit struct {
	Which  int
	Enable bool
	RLim   syscall.Rlimit
}

type ITimerVal struct {
	ItInterval TimeVal
	ItValue    TimeVal
}

type TimeVal struct {
	TvSec  uint64
	TvUsec uint64
}

type ExecRLimit struct {
	TimeLimit     int // 时间限制 (ms)
	RealTimeLimit int // 真实时间限制 (ms, 触发SIGALRM)
	MemoryLimit   int // 内存限制 (KB)
	FileSizeLimit int // 文件读写限制 (B)
	StackLimit    int // 栈大小限制 (KB，0表示用内存限制的值，-1表示不限制，建议设置为2倍。Mac下有坑，不要去设置。)
}

type RlimitOptions struct {
	Rlimits     []RLimit
	ITimerValue ITimerVal
}

// 解析ExecRLimit结构体并获取setrlimit操作需要的信息
func GetRlimitOptions(sysRlimit *ExecRLimit) *RlimitOptions {
	// Make stack limit
	stackLimit := uint64(sysRlimit.StackLimit)
	if stackLimit <= 0 {
		stackLimit = uint64(sysRlimit.MemoryLimit * 2)
	}

	return &RlimitOptions{
		Rlimits: []RLimit{
			// Set time limit: RLIMIT_CPU
			{
				Which:  syscall.RLIMIT_CPU,
				Enable: sysRlimit.TimeLimit > 0,
				RLim: syscall.Rlimit{
					Cur: uint64(math.Ceil(float64(sysRlimit.TimeLimit) / 1000.0)),
					Max: uint64(math.Ceil(float64(sysRlimit.TimeLimit) / 1000.0)),
				},
			},
			// Set memory limit: RLIMIT_DATA
			{
				Which:  syscall.RLIMIT_DATA,
				Enable: sysRlimit.MemoryLimit > 0,
				RLim: syscall.Rlimit{
					Cur: uint64(sysRlimit.MemoryLimit * 1024),
					Max: uint64(sysRlimit.MemoryLimit * 1024 * 2),
				},
			},
			// Set memory limit: RLIMIT_AS
			{
				Which:  syscall.RLIMIT_AS,
				Enable: sysRlimit.MemoryLimit > 0,
				RLim: syscall.Rlimit{
					Cur: uint64(sysRlimit.MemoryLimit * 1024 * 2),
					Max: uint64(sysRlimit.MemoryLimit*1024*2 + 1024),
				},
			},
			// Set stack limit (坑：macos不要去搞这个!)
			{
				Which:  syscall.RLIMIT_STACK,
				Enable: stackLimit > 0 && sysRlimit.StackLimit >= 0 && runtime.GOOS != "darwin",
				RLim: syscall.Rlimit{
					Cur: stackLimit * 1024,
					Max: stackLimit*1024 + 1024,
				},
			},
			// Set file size limit: RLIMIT_FSIZE
			{
				Which:  syscall.RLIMIT_FSIZE,
				Enable: sysRlimit.FileSizeLimit > 0,
				RLim: syscall.Rlimit{
					Cur: uint64(sysRlimit.FileSizeLimit),
					Max: uint64(sysRlimit.FileSizeLimit),
				},
			},
		},
		ITimerValue: ITimerVal{
			ItInterval: TimeVal{
				TvSec:  uint64(math.Floor(float64(sysRlimit.RealTimeLimit) / 1000.0)),
				TvUsec: uint64(sysRlimit.RealTimeLimit % 1000 * 1000),
			},
			ItValue: TimeVal{
				TvSec:  uint64(math.Floor(float64(sysRlimit.RealTimeLimit) / 1000.0)),
				TvUsec: uint64(sysRlimit.RealTimeLimit % 1000 * 1000),
			},
		},
	}
}
