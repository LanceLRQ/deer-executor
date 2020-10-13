package executor

import (
	"fmt"
	"path"
	"syscall"
)

// 分析进程退出状态
func (session *JudgeSession) analysisExitStatus(rst *JudgeResult, pinfo *ProcessInfo, specialJudge bool) error {
	ru := pinfo.Rusage
	status := pinfo.Status

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
	pinfo, err := session.runProgram()
	if err != nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = err.Error()
		return err
	}

	err = session.analysisExitStatus(judgeResult, pinfo, false)
	if err != nil {
		return err
	}
	return nil
}

// 执行评测
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
	session.SpecialJudge.LogFile = path.Join(session.SessionDir, "judger.log")

	// 运行judge程序
	err = session.judge(&judgeResult)
	if err != nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = err.Error()
		return judgeResult, err
	}

	return judgeResult, nil
}
