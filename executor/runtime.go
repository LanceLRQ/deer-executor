package executor

import (
	"fmt"
	"github.com/docker/docker/pkg/reexec"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)


// 运行目标进程
func (options *JudgeOptions)runTargetProgram() (*exec.Cmd, error) {

	params := []string {
		"targetProgram",
		strconv.Itoa(options.TimeLimit),
		strconv.Itoa(options.MemoryLimit),
		strconv.Itoa(options.RealTimeLimit),
		strconv.Itoa(options.Uid),
	}
	params = append(params, options.Commands...)

	target := reexec.Command(params...)

	tcin, err := OpenFile(options.TestCaseIn, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil { return target, err }
	target.Stdin = tcin

	pout, err := os.OpenFile(options.ProgramOut, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil { return target, err }
	target.Stdout = pout

	perr, err := os.OpenFile(options.ProgramError, os.O_WRONLY | os.O_CREATE, 0644)
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


func (options *JudgeOptions) analysisExitStatus(rst *JudgeResult, cmd *exec.Cmd, specialJudge bool) error {
	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	ru := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	rst.TimeUsed = int(ru.Utime.Sec * 1000 + int64(ru.Utime.Usec) / 1000 + ru.Stime.Sec * 1000 + int64(ru.Stime.Usec) / 1000)
	rst.MemoryUsed = int(ru.Minflt * int64(syscall.Getpagesize() / 1024 ))

	// If process stopped with a signal
	if status.Signaled() {
		sig := status.Signal()
		rst.ReSignum = int(sig)
		if sig == syscall.SIGSEGV {
			// MLE or RE can also get SIGSEGV signal.
			if rst.MemoryUsed > options.MemoryLimit {
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
			if rst.TimeUsed > (options.TimeLimit - 100) {
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
		if rst.TimeUsed > options.TimeLimit {
			rst.JudgeResult = JudgeFlagMLE
		} else if rst.MemoryUsed > options.MemoryLimit {
			rst.JudgeResult = JudgeFlagMLE
		} else {
			rst.JudgeResult = JudgeFlagAC
		}
	}
	return nil
}

// 基于JudgeOptions进行评测调度
func (options *JudgeOptions) judge(judgeResult *JudgeResult) error {
	target, err := options.runTargetProgram()
	if target == nil && err != nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = err.Error()
		return err
	} else {
		err = options.analysisExitStatus(judgeResult, target, false)
		if err != nil {
			judgeResult.JudgeResult = JudgeFlagSE
			judgeResult.SeInfo = err.Error()
			return err
		}
	}
	return nil
}

func (options *JudgeOptions)RunJudge() (JudgeResult, error) {
	judgeResult := JudgeResult{}
	// 获取对应的编译器提供程序
	compiler, err := options.getCompiler("")
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
	options.Commands = compiler.GetRunArgs()

	// 清理输出文件，以免文件数据错误
	_ = os.Remove(options.ProgramOut)
	_ = os.Remove(options.ProgramError)
	_ = os.Remove(options.SpecialJudge.Stdout)
	_ = os.Remove(options.SpecialJudge.Stderr)
	_ = os.Remove(options.SpecialJudge.Logfile)

	// 运行judge程序
	err = options.judge(&judgeResult)
	if err != nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = err.Error()
		return judgeResult, err
	}

	return judgeResult, nil
}



// 目标程序子进程
func RunTargetProgramProcess() {
	if len(os.Args) < 5 {
		log.Fatal("params error.")
		return
	}
	timeLimit, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		log.Fatal("parse time limit number error.")
		return
	}
	memoryLimit, err := strconv.ParseInt(os.Args[2], 10, 32)
	if err != nil {
		log.Fatal("parse memory limit number error.")
		return
	}
	realTimeLimit, err := strconv.ParseInt(os.Args[3], 10, 32)
	if err != nil {
		log.Fatal("parse real time limit number error.")
		return
	}
	uid, err := strconv.ParseInt(os.Args[4], 10, 32)
	if err != nil {
		log.Fatal("parse uid number error.")
		return
	}
	commands := os.Args[5:len(os.Args)]
	// Set UID
	if uid > -1 {
		err := syscall.Setuid(int(uid))
		if err != nil {
			log.Fatalf("set resource limit error: %s", err.Error())
			return
		}
	}
	// Set Resource Limit
	err = setLimit(int(timeLimit), int(memoryLimit), int(realTimeLimit))
	if err != nil {
		log.Fatalf("set resource limit error: %s", err.Error())
		return
	}
	// Run Program
	if len(commands) > 1 {
		_ = syscall.Exec(commands[0], commands[1:], CommonEnvs)
	} else {
		_ = syscall.Exec(commands[0], nil, CommonEnvs)
	}
}

// 特判程序子进程
func RunSpecialJudgeProgramProcess() {

}
