package executor

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
)


// 运行评测进程
func (session *JudgeSession)runProgramCommon(judger bool, pipeMode bool, pipeStd []int) (*ProcessInfo, error) {
	pinfo := ProcessInfo{}
	pid, fds, err := session.runProgramProcess(judger, pipeMode, pipeStd)
	if err != nil {
		return nil, err
	}
	pinfo.Pid = pid
	// before wait4, do something~

	// Wait4
	_, err = syscall.Wait4(int(pid), &pinfo.Status, syscall.WUNTRACED, &pinfo.Rusage)
	if err != nil {
		return nil, err
	}

	// Close Files
	for _, fd := range fds {
		if fd > 0 {
			_ = syscall.Close(fd)
		}
	}

	return &pinfo, err
}

// 运行目标程序
func (session *JudgeSession)runNormalJudge() (*ProcessInfo, error) {
	return session.runProgramCommon(false, false, nil)
}

// 运行特殊评测
func (session *JudgeSession)runSpecialJudge() (*ProcessInfo, *ProcessInfo, error) {
	if session.SpecialJudge.Mode == SpecialJudgeModeChecker {
		targetInfo, err := session.runProgramCommon(false, false, nil)
		judgerInfo, err := session.runProgramCommon(true, false, nil)
		return targetInfo, judgerInfo, err
	} else if session.SpecialJudge.Mode == SpecialJudgeModeInteractive {

		fdjudger, err := getPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("create pipe error: %s", err.Error())
		}

		fdtarget, err := getPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("create pipe error: %s", err.Error())
		}

		targetInfo, err := session.runProgramCommon(false, true, []int{fdtarget[0], fdjudger[1]})
		judgerInfo, err := session.runProgramCommon(true, true, []int{fdjudger[0], fdtarget[1]})
		return targetInfo, judgerInfo, err
	}
	return nil, nil, fmt.Errorf("unkonw special judge mode")
}


// 目标程序子进程
func (session *JudgeSession)runProgramProcess(judger bool, pipeMode bool, pipeStd []int) (uintptr, []int, error) {
	var (
		logfile *os.File
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
		var logWriter *bufio.Writer
		if judger {
			logfile, err = os.OpenFile(session.SpecialJudge.LogFile, os.O_WRONLY|os.O_CREATE, 0644)
		} else {
			logfile, err = os.OpenFile(session.ProgramLog, os.O_WRONLY|os.O_CREATE, 0644)
		}
		if err != nil {
			panic("cannot create program.log")
			return 0, fds, err
		} else {
			logWriter = bufio.NewWriter(logfile)
		}

		if pipeMode {
			// Direct Program's Pipe[Read] to Stdin
			err = syscall.Dup2(pipeStd[0], syscall.Stdin)
			if err != nil {
				return 0, fds, err
			}
			// Direct Judger's Pipe[Write] to Stdout
			err = syscall.Dup2(pipeStd[1], syscall.Stdout)
			if err != nil {
				return 0, fds, err
			}
		} else {
			// Redirect testCaseIn to STDIN
			if judger {
				if session.SpecialJudge.RedirectProgramOut {
					fds[0], err = redirectFileDescriptor(syscall.Stdin, session.ProgramOut, os.O_RDONLY, 0)
				} else {
					fds[0], err = redirectFileDescriptor(syscall.Stdin, session.TestCaseIn, os.O_RDONLY, 0)
				}
			} else {
				fds[0], err = redirectFileDescriptor(syscall.Stdin, session.TestCaseIn, os.O_RDONLY, 0)
			}
			if err != nil {
				_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]direct stdin error: %s\n", err.Error()))
				return 0, fds, err
			}

			// Redirect userOut to STDOUT
			if judger {
				fds[1], err = redirectFileDescriptor(syscall.Stdout, session.SpecialJudge.Stdout, os.O_WRONLY|os.O_CREATE, 0644)
			} else {
				fds[1], err = redirectFileDescriptor(syscall.Stdout, session.ProgramOut, os.O_WRONLY|os.O_CREATE, 0644)
			}
			if err != nil {
				_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]direct stdout error: %s\n", err.Error()))
				return 0, fds, err
			}
		}

		// Redirect programError to STDERR
		if judger {
			fds[2], err = redirectFileDescriptor(syscall.Stderr, session.SpecialJudge.Stderr, os.O_WRONLY|os.O_CREATE, 0644)
		} else {
			fds[2], err = redirectFileDescriptor(syscall.Stderr, session.ProgramError, os.O_WRONLY|os.O_CREATE, 0644)
		}
		if err != nil {
			_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]direct stderr error: %s\n", err.Error()))
			return 0, fds, err
		}

		// Set UID
		if session.Uid > -1 {
			err = syscall.Setuid(session.Uid)
			if err != nil {
				_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]set resource limit error: %s\n", err.Error()))
				return 0, fds, err
			}
		}

		// Set Resource Limit
		if judger {
			err = setLimit(session.SpecialJudge.TimeLimit, session.SpecialJudge.MemoryLimit, session.RealTimeLimit)
		} else {
			err = setLimit(session.TimeLimit, session.MemoryLimit, session.RealTimeLimit)
		}
		if err != nil {
			_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]set resource limit error: %s", err.Error()))
			return 0, fds, err
		}

		if judger {
			// Run Judger (Testlib compatible)
			// ./checker <input-file> <output-file> <answer-file> <report-file>
			args := []string{
				session.SpecialJudge.Checker,
				session.TestCaseIn,
				session.TestCaseOut,
				session.ProgramOut,
				session.SpecialJudge.ReportFile,
			}
			err = syscall.Exec(session.SpecialJudge.Checker, args, nil)
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
	}
	// parent process
	return pid, fds, nil
}
