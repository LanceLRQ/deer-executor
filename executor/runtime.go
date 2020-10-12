package executor

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"syscall"
)


// 普通评测进程
func (session *JudgeSession)runProgramNormal() (*syscall.WaitStatus, *syscall.Rusage, error) {
	pid, fds, err := session.runTargetProgramProcess()
	if err != nil {
		return nil, nil, err
	}
	// before wait4, do something~
	var (
		status syscall.WaitStatus
		ru syscall.Rusage
	)
	// Wait4
	_, err = syscall.Wait4(int(pid), &status, syscall.WUNTRACED, &ru)
	if err != nil {
		return nil, nil, err
	}

	// Close Files
	for _, fd := range fds {
		_ = syscall.Close(fd)
	}

	return &status, &ru, err
}


func (session *JudgeSession) analysisExitStatus(rst *JudgeResult, status *syscall.WaitStatus, ru *syscall.Rusage, specialJudge bool) error {
	rst.TimeUsed = int(ru.Utime.Sec * 1000 + int64(ru.Utime.Usec) / 1000 + ru.Stime.Sec * 1000 + int64(ru.Stime.Usec) / 1000)
	rst.MemoryUsed = int(ru.Minflt * int64(syscall.Getpagesize() / 1024 ))

	// If process stopped with a signal
	if status.Signaled() {
		sig := status.Signal()
		rst.ReSignum = int(sig)
		if sig == syscall.SIGSEGV {
			// MLE or RE can also get SIGSEGV signal.
			if rst.MemoryUsed > session.MemoryLimit {
				rst.JudgeResult = JudgeFlagMLE
			} else {
				rst.JudgeResult = JudgeFlagRE
			}
		} else if sig == syscall.SIGXFSZ {
			// SIGXFSZ signal means OLE
			rst.JudgeResult = JudgeFlagOLE
		} else if sig == syscall.SIGALRM || sig == syscall.SIGVTALRM || sig == syscall.SIGXCPU {
			// Normal TLE signal
			rst.JudgeResult = JudgeFlagTLE
		} else if sig == syscall.SIGKILL {
			// Sometimes MLE might get SIGKILL signal.
			// So if real time used lower than TIME_LIMIT - 100, it might be a TLE error.
			if rst.TimeUsed > (session.TimeLimit - 100) {
				rst.JudgeResult = JudgeFlagTLE
			} else {
				rst.JudgeResult = JudgeFlagMLE
			}
		} else {
			// Otherwise, called runtime error.
			rst.JudgeResult = JudgeFlagRE
		}
	} else {
		// Sometimes setrlimit doesn't work accurately.
		if rst.TimeUsed > session.TimeLimit {
			rst.JudgeResult = JudgeFlagMLE
		} else if rst.MemoryUsed > session.MemoryLimit {
			rst.JudgeResult = JudgeFlagMLE
		} else {
			rst.JudgeResult = JudgeFlagAC
		}
	}
	return nil
}

// 基于JudgeOptions进行评测调度
func (session *JudgeSession) judge(judgeResult *JudgeResult) error {
	status, ru, err := session.runProgramNormal()
	if err != nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = err.Error()
		return err
	}

	err = session.analysisExitStatus(judgeResult, status, ru, false)
	if err != nil {
		return err
	}
	return nil
}

func (session *JudgeSession)RunJudge() (JudgeResult, error) {
	judgeResult := JudgeResult{}
	// 获取对应的编译器提供程序
	compiler, err := session.getCompiler("")
	if err != nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = err.Error()
		return judgeResult, err
	}
	// 编译程序
	success, ceinfo := compiler.Compile()
	if !success {
		judgeResult.JudgeResult = JudgeFlagCE
		judgeResult.CeInfo = ceinfo
		return judgeResult, fmt.Errorf("compile error:\n%s", ceinfo)
	}
	// 清理工作目录
	defer compiler.Clean()
	// 获取执行指令
	session.Commands = compiler.GetRunArgs()

	// 创建相关的文件路径
	session.ProgramOut = path.Join(session.SessionDir, "program.out")
	session.ProgramError = path.Join(session.SessionDir, "program.err")
	session.ProgramLog = path.Join(session.SessionDir, "program.log")
	session.SpecialJudge.Stdout = path.Join(session.SessionDir, "judger.out")
	session.SpecialJudge.Stderr = path.Join(session.SessionDir, "judger.err")
	session.SpecialJudge.Logfile = path.Join(session.SessionDir, "judger.log")

	// 运行judge程序
	err = session.judge(&judgeResult)
	if err != nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = err.Error()
		return judgeResult, err
	}

	return judgeResult, nil
}



// 目标程序子进程
func (session *JudgeSession)runTargetProgramProcess() (uintptr, []int, error) {
	var (
		err error
		pid uintptr
		fds []int
	)

	fds = make([]int, 3)

	// Fork a new process
	pid, err = forkProc()
	if err != nil {
		return 0, nil, fmt.Errorf("fork process error: %s", err.Error())
	}

	if pid == 0 {
		var logWriter *bufio.Writer
		logfile, err := os.OpenFile(session.ProgramLog, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic("cannot create program.log")
			return 0, nil, err
		} else {
			logWriter = bufio.NewWriter(logfile)
		}

		// Redirect testCaseIn to STDIN
		fds[0], err = redirectFileDescriptor(syscall.Stdin, session.TestCaseIn, os.O_RDONLY, 0)
		if err != nil {
			_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]direct stdin error: %s\n", err.Error()))
			return 0, nil, err
		}

		// Redirect userOut to STDOUT
		fds[1], err = redirectFileDescriptor(syscall.Stdout, session.ProgramOut, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]direct stdout error: %s\n", err.Error()))
			return 0, nil, err
		}

		// Redirect programError to STDERR
		fds[2], err = redirectFileDescriptor(syscall.Stderr, session.ProgramError, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]direct stderr error: %s\n", err.Error()))
			return 0, nil, err
		}

		// Set UID
		if session.Uid > -1 {
			err = syscall.Setuid(session.Uid)
			if err != nil {
				_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]set resource limit error: %s\n", err.Error()))
				return 0, nil, err
			}
		}

		// Set Resource Limit
		err = setLimit(session.TimeLimit, session.MemoryLimit, session.RealTimeLimit)
		if err != nil {
			_, _ = logWriter.WriteString(fmt.Sprintf("[system_error]set resource limit error: %s", err.Error()))
			return 0, nil, err
		}

		// Run Program
		commands := session.Commands
		if len(commands) > 1 {
			_ = syscall.Exec(commands[0], commands[1:], CommonEnvs)
		} else {
			_ = syscall.Exec(commands[0], nil, CommonEnvs)
		}
		// it won't be run.
	}
	// parent process
	return pid, fds, nil
}

// 特判程序子进程
func RunSpecialJudgeProgramProcess() {

}
