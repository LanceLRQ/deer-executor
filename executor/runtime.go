package executor

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"syscall"
)


// 运行评测进程
func (session *JudgeSession)runProgramCommon(rst *TestCaseResult, judger bool, pipeMode bool, pipeStd []int) (*ProcessInfo, error) {
	pinfo := ProcessInfo{}
	pid, fds, err := runProgramProcess(session, rst, judger, pipeMode, pipeStd)
	if err != nil {
		if pid == 0 {
			// 如果是子进程错误了，输出到程序的error去
			panic(err)
		}
		return nil, err
	}
	pinfo.Pid = pid
	// before wait4, do something~

	// Wait4
	_, err = syscall.Wait4(int(pid), &pinfo.Status, syscall.WUNTRACED, &pinfo.Rusage)
	if err != nil {
		return nil, err
	}

	if !pipeMode {
		// Close Files
		for _, fd := range fds {
			if fd > 0 {
				_ = syscall.Close(fd)
			}
		}
	}

	return &pinfo, err
}

// 运行交互评测进程
func (session *JudgeSession)runProgramAsync(rst *TestCaseResult, judger bool, pipeMode bool, pipeStd []int, info chan *ProcessInfo) error {
	tpid, fds, err := runProgramProcess(session, rst, judger, pipeMode, pipeStd)
	if err != nil {
		if tpid == 0 {
			// 如果是子进程错误了(没能正确执行到目标程序里)，输出到程序的error去
			panic(err)
		}
		return err
	}

	go func(pid uintptr) {
		pinfo := ProcessInfo{}
		pinfo.Pid = pid
		// Wait4
		_, err = syscall.Wait4(int(pid), &pinfo.Status, syscall.WUNTRACED, &pinfo.Rusage)
		if err != nil {
			info <- &pinfo
			return
		}

		// Close Files
		if !pipeMode {
			for _, fd := range fds {
				if fd > 0 {
					_ = syscall.Close(fd)
				}
			}
		}
		info <- &pinfo
	}(tpid)

	return nil
}


// 运行目标程序
func (session *JudgeSession)runNormalJudge(rst *TestCaseResult) (*ProcessInfo, error) {
	return session.runProgramCommon(rst, false, false, nil)
}

// 运行特殊评测
func (session *JudgeSession)runSpecialJudge(rst *TestCaseResult) (*ProcessInfo, *ProcessInfo, error) {
	if session.JudgeConfig.SpecialJudge.Mode == SpecialJudgeModeChecker {
		targetInfo, err := session.runProgramCommon(rst, false, false, nil)
		if err != nil {
			return targetInfo, nil, err
		}
		judgerInfo, err := session.runProgramCommon(rst, true, false, nil)
		return targetInfo, judgerInfo, err
	} else if session.JudgeConfig.SpecialJudge.Mode == SpecialJudgeModeInteractive {

		fdjudger, err := getPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("create pipe error: %s", err.Error())
		}

		fdtarget, err := getPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("create pipe error: %s", err.Error())
		}

		targetInfoChan, judgerInfoChan := make(chan *ProcessInfo), make(chan *ProcessInfo)
		var targetInfo, judgerInfo *ProcessInfo

		err = session.runProgramAsync(rst, false, true, []int {fdtarget[0], fdjudger[1]}, targetInfoChan)
		if err != nil {
			return nil, nil, err
		}
		err = session.runProgramAsync(rst, true, true, []int {fdjudger[0], fdtarget[1]}, judgerInfoChan)
		if err != nil {
			return nil, nil, err
		}
		targetInfo = <- targetInfoChan
		judgerInfo = <- judgerInfoChan
		return targetInfo, judgerInfo, err
	}
	return nil, nil, fmt.Errorf("unkonw special judge mode")
}

func getSpecialJudgerPath(session *JudgeSession, rst *TestCaseResult) []string {
	tci, err := filepath.Abs(path.Join(session.ConfigDir, rst.TestCaseIn))
	if err != nil {
		tci = path.Join(session.ConfigDir, rst.TestCaseIn)
	}
	tco, err := filepath.Abs(path.Join(session.ConfigDir, rst.TestCaseOut))
	if err != nil {
		tco = path.Join(session.ConfigDir, rst.TestCaseOut)
	}
	po, err := filepath.Abs(path.Join(session.SessionDir, rst.ProgramOut))
	if err != nil {
		po = path.Join(session.SessionDir, rst.ProgramOut)
	}
	jr, err := filepath.Abs(path.Join(session.SessionDir, rst.JudgerReport))
	if err != nil {
		jr = path.Join(session.SessionDir, rst.JudgerReport)
	}
	args := []string{
		tci,
		tco,
		po,
		jr,
	}
	return args
}


func getLimitation(session *JudgeSession) (int, int, int, int, int) {
	langName := session.Compiler.GetName()
	memoryLimitExtend := 0
	jitMem, ok := MemorySizeForJIT[langName]
	if ok {
		memoryLimitExtend = jitMem
	}
	limitation, ok := session.JudgeConfig.Limitation[langName]
	if ok {
		return limitation.TimeLimit,
			limitation.MemoryLimit + memoryLimitExtend,
			limitation.RealTimeLimit,
			limitation.FileSizeLimit,
			memoryLimitExtend
	}
	return session.JudgeConfig.TimeLimit,
		session.JudgeConfig.MemoryLimit + memoryLimitExtend,
		session.JudgeConfig.RealTimeLimit,
		session.JudgeConfig.FileSizeLimit,
		memoryLimitExtend
}

// 目标程序子进程
func runProgramProcess(session *JudgeSession, rst *TestCaseResult, judger bool, pipeMode bool, pipeStd []int) (uintptr, []int, error) {
	var (
		err error
		pid uintptr
		fds []int
	)

	fds = make([]int, 3)

	// Fork a new process
	pid, err = forkProc()
	if err != nil {
		return 0, fds, fmt.Errorf("fork process error: %s", err.Error())
	}

	if pid == 0 {
		if pipeMode {
			// Direct Pipe[Read] to Stdin
			err = syscall.Dup2(pipeStd[0], syscall.Stdin)
			if err != nil {
				return 0, fds, err
			}
			// Direct Pipe[Write] to Stdout
			err = syscall.Dup2(pipeStd[1], syscall.Stdout)
			if err != nil {
				return 0, fds, err
			}
		} else {
			// Redirect testCaseIn to STDIN
			if judger {
				if session.JudgeConfig.SpecialJudge.RedirectProgramOut {
					fds[0], err = redirectFileDescriptor(
						syscall.Stdout,
						path.Join(session.SessionDir, rst.ProgramOut),
						os.O_RDONLY,
						0,
					)
				} else {
					fds[0], err = redirectFileDescriptor(
						syscall.Stdin,
						path.Join(session.ConfigDir, rst.TestCaseIn),
						os.O_RDONLY,
						0,
					)
				}
			} else {
				fds[0], err = redirectFileDescriptor(
					syscall.Stdin,
					path.Join(session.ConfigDir, rst.TestCaseIn),
					os.O_RDONLY,
					0,
				)
			}
			if err != nil {
				return 0, fds, err
			}

			// Redirect userOut to STDOUT
			if judger {
				fds[1], err = redirectFileDescriptor(
					syscall.Stdout,
					path.Join(session.SessionDir, rst.JudgerOut),
					os.O_WRONLY|os.O_CREATE, 0644,
				)
			} else {
				fds[1], err = redirectFileDescriptor(
					syscall.Stdout,
					path.Join(session.SessionDir, rst.ProgramOut),
					os.O_WRONLY|os.O_CREATE,
					0644,
				)
			}
			if err != nil {
				return 0, fds, err
			}
		}

		// Redirect programError to STDERR
		if judger {
			fds[2], err = redirectFileDescriptor(
				syscall.Stderr,
				path.Join(session.SessionDir, rst.JudgerError),
				os.O_WRONLY|os.O_CREATE,
				0644,
			)
		} else {
			fds[2], err = redirectFileDescriptor(
				syscall.Stderr,
				path.Join(session.SessionDir, rst.ProgramError),
				os.O_WRONLY|os.O_CREATE,
				0644,
			)
		}
		if err != nil {
			return 0, fds, err
		}

		// Set UID
		if session.JudgeConfig.Uid > -1 {
			err = syscall.Setuid(session.JudgeConfig.Uid)
			if err != nil {
				return 0, fds, err
			}
		}

		// Set Resource Limit
		if judger {
			err = setLimit(
				session.JudgeConfig.SpecialJudge.TimeLimit,
				session.JudgeConfig.SpecialJudge.MemoryLimit,
				session.JudgeConfig.RealTimeLimit,
				session.JudgeConfig.FileSizeLimit,
				)
		} else {
			tl, ml, rtl, fsl, _ := getLimitation(session)
			err = setLimit(tl, ml, rtl, fsl )
		}
		if err != nil {
			return 0, fds, err
		}

		if judger {
			// Run Judger (Testlib compatible)
			// ./checker <input-file> <output-file> <answer-file> <report-file>
			args := getSpecialJudgerPath(session, rst)
			_ = syscall.Exec(session.JudgeConfig.SpecialJudge.Checker, args, nil)
		} else {
			// Run Program
			commands := session.Commands
			if len(commands) > 1 {
				_ = syscall.Exec(commands[0], commands[1:], CommonEnvs)
			} else {
				_ = syscall.Exec(commands[0], nil, CommonEnvs)
			}
		}
		// it won't be run.
	} else if pid < 0 {
		return 0, fds, fmt.Errorf("fork process error: pid < 0")
	}
	// parent process
	return pid, fds, nil
}
