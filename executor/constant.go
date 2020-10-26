package executor

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	JudgeFlagAC 	 						= 0   					// 0 Accepted
	JudgeFlagPE 	 						= 1	    				// 1 Presentation Error
	JudgeFlagTLE 							= 2						// 2 Time Limit Exceeded
	JudgeFlagMLE 	 						= 3						// 3 Memory Limit Exceeded
	JudgeFlagWA 	 						= 4	    				// 4 Wrong Answer
	JudgeFlagRE 	 						= 5	    				// 5 Runtime Error
	JudgeFlagOLE 	 						= 6						// 6 Output Limit Exceeded
	JudgeFlagCE 	 						= 7	    				// 7 Compile Error
	JudgeFlagSE 						 	= 8     				// 8 System Error

	JudgeFlagSpecialJudgeTimeout 		 	= 10    				// 10 Special Judger Time OUT
	JudgeFlagSpecialJudgeError 			 	= 11    				// 11 Special Judger ERROR
	JudgeFlagSpecialJudgeRequireChecker  	= 12 					// 12 Special Judger Finish, Need Standard Checkup
)

const (
	SpecialJudgeModeDisabled 				= 0
	SpecialJudgeModeChecker 				= 1
	SpecialJudgeModeInteractive 			= 2

	SpecialJudgeTimeLimit 					= 1 * 1000				// Unit: ms
	SpecialJudgeMemoryLimit 				= 256 * 1024			// Unit: kb
)

var SignalNumberMap = map[int][]string {
	1: []string{"SIGHUP", "Hangup (POSIX)."},
	2:  []string{"SIGINT", "Interrupt (ANSI)."},
	3:  []string{"SIGQUIT", "Quit (POSIX)."},
	4:  []string{"SIGILL", "Illegal instruction (ANSI)."},
	5:  []string{"SIGTRAP", "Trace trap (POSIX)."},
	6:  []string{"SIGABRT", "Abort (ANSI)."},
	//6:  []string{"SIGIOT", "IOT trap (4.2 BSD)."},
	7:  []string{"SIGBUS", "BUS error (4.2 BSD)."},
	8:  []string{"SIGFPE", "Floating-point exception (ANSI)."},
	9:  []string{"SIGKILL", "Kill, unblockable (POSIX)."},
	10:  []string{"SIGUSR1", "User-defined signal 1 (POSIX)."},
	11:  []string{"SIGSEGV", "Segmentation violation (ANSI)."},
	12:  []string{"SIGUSR2", "User-defined signal 2 (POSIX)."},
	13:  []string{"SIGPIPE", "Broken pipe (POSIX)."},
	14:  []string{"SIGALRM", "Alarm clock (POSIX)."},
	15:  []string{"SIGTERM", "Termination (ANSI)."},
	16:  []string{"SIGSTKFLT", "Stack fault."},
	17:  []string{"SIGCHLD", "Child status has changed (POSIX)."},
	18:  []string{"SIGCONT", "Continue (POSIX)."},
	19:  []string{"SIGSTOP", "Stop, unblockable (POSIX)."},
	20:  []string{"SIGTSTP", "Keyboard stop (POSIX)."},
	21:  []string{"SIGTTIN", "Background read from tty (POSIX)."},
	22:  []string{"SIGTTOU", "Background write to tty (POSIX)."},
	23:  []string{"SIGURG", "Urgent condition on socket (4.2 BSD)."},
	24:  []string{"SIGXCPU", "CPU limit exceeded (4.2 BSD)."},
	25:  []string{"SIGXFSZ", "File size limit exceeded (4.2 BSD)."},
	26:  []string{"SIGVTALRM", "Virtual alarm clock (4.2 BSD)."},
	27:  []string{"SIGPROF", "Profiling alarm clock (4.2 BSD)."},
	28:  []string{"SIGWINCH", "Window size change (4.3 BSD, Sun)."},
	29:  []string{"SIGIO", "I/O now possible (4.2 BSD)."},
	30:  []string{"SIGPWR", "Power failure restart (System V)."},
	31:  []string{"SIGSYS", "Bad system call."},
}

var FlagMeansMap = map[int]string {
	0: "Accepted",
	1: "Presentation Error",
	2: "Time Limit Exceeded",
	3: "Memory Limit Exceeded",
	4: "Wrong Answer",
	5: "Runtime Error",
	6: "Output Limit Exceeded",
	7: "Compile Error",
	8: "System Error",
	9: "Special Judger Time OUT",
	10: "Special Judger ERROR",
	11: "Special Judger Finish, Need Standard Checkup",
}

// 给动态语言、带虚拟机的语言设定虚拟机自身的初始内存大小
var MemorySizeForJIT = map[string]int {
	"gcc": 			0,
	"g++": 			0,
	"java": 		393216,			// java
	"python2": 		65536,			// py2
	"python3": 		65536,			// py3
	"nodejs": 		262144,			// js
	"golang":		0,
	"php": 			131072,			// php
	"ruby": 		65536,			// ruby
}

func PlaceMemorySizeForJIT(configFile string) error {
	if configFile != "" {
		_, err := os.Stat(configFile)
		// ignore
		if os.IsNotExist(err) {
			return nil
		}
		cbody, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(cbody, &MemorySizeForJIT)
		if err != nil {
			return err
		}
	}
	return nil
}