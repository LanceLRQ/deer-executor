package executor

import (
	"bufio"
	"fmt"
	"github.com/docker/docker/pkg/reexec"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)


// 运行目标进程
func (session *JudgeSession)runTargetProgram() (*exec.Cmd, error) {

	opt := ObjectToJSONString(session)
	target := reexec.Command("targetProgram", opt)

	tcin, err := OpenFile(session.TestCaseIn, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil { return target, err }
	target.Stdin = tcin

	pout, err := os.OpenFile(session.ProgramOut, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil { return target, err }
	target.Stdout = pout

	perr, err := os.OpenFile(session.ProgramError, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil { return target, err }
	target.Stderr = perr

	if err := target.Start(); err != nil {
		return target, fmt.Errorf("failed to run target program: %s", err)
	}
	if err := target.Wait(); err != nil {
		return target, fmt.Errorf("failed to wait command: %s", err)
	}
	return target, nil
}


func (session *JudgeSession) analysisExitStatus(rst *JudgeResult, cmd *exec.Cmd, specialJudge bool) error {
	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	ru := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	rst.TimeUsed = int(ru.Utime.Sec * 1000 + int64(ru.Utime.Usec) / 1000 + ru.Stime.Sec * 1000 + int64(ru.Stime.Usec) / 1000)
	rst.MemoryUsed = int(ru.Minflt * int64(syscall.Getpagesize() / 1024 ))

	// Fix time used & mem used
	if sysLog, err := os.OpenFile(path.Join(session.SessionDir, "sys.log"), syscall.O_RDONLY | syscall.O_NONBLOCK, 0644); err == nil {
		reader := bufio.NewReader(sysLog)

		if line, _, err := reader.ReadLine(); err == nil {
			if tu, err := strconv.ParseInt(string(line), 10, 32); err == nil {
				rst.TimeUsed -= int(tu)
			}
		}
		if line, _, err := reader.ReadLine(); err == nil {
			if mu, err := strconv.ParseInt(string(line), 10, 32); err == nil {
				rst.MemoryUsed -= int(mu)
			}
		}
	}

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
	target, err := session.runTargetProgram()
	if err != nil && target == nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = err.Error()
		return err
	} else {
		err = session.analysisExitStatus(judgeResult, target, false)
		if err != nil {
			judgeResult.JudgeResult = JudgeFlagSE
			judgeResult.SeInfo = err.Error()
			return err
		}
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
func RunTargetProgramProcess() {
	var ru syscall.Rusage
	payload := os.Args[1]
	session := JudgeSession{}

	if !JSONStringObject(payload, &session) {
		log.Fatal("[system_error]parse judge session error")
		return
	}
	// Set UID
	if session.Uid > -1 {
		err := syscall.Setuid(session.Uid)
		if err != nil {
			log.Fatalf("[system_error]set resource limit error: %s", err.Error())
			return
		}
	}

	_ = syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
	tu := int(ru.Utime.Sec * 1000 + int64(ru.Utime.Usec) / 1000 + ru.Stime.Sec * 1000 + int64(ru.Stime.Usec) / 1000)
	mu := int(ru.Minflt * int64(syscall.Getpagesize() / 1024 ))

	// Set Resource Limit
	err := setLimit(session.TimeLimit + tu, session.MemoryLimit + mu, session.RealTimeLimit)
	if err != nil {
		log.Fatalf("[system_error]set resource limit error: %s", err.Error())
		return
	}

	// Save current rusage of judger
	if sysLog, err := os.OpenFile(path.Join(session.SessionDir, "sys.log"), syscall.O_WRONLY | syscall.O_CREAT, 0644); err == nil {


		_, _ = sysLog.WriteString(strconv.Itoa(tu) + "\n")
		_, _ = sysLog.WriteString(strconv.Itoa(mu) + "\n")
		_ = sysLog.Close()
	}

	// Run Program
	commands := session.Commands
	if len(commands) > 1 {
		_ = syscall.Exec(commands[0], commands[1:], CommonEnvs)
	} else {
		_ = syscall.Exec(commands[0], nil, CommonEnvs)
	}
}

// 特判程序子进程
func RunSpecialJudgeProgramProcess() {

}
