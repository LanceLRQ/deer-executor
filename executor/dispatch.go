package executor

import (
	"fmt"
	"path"
	"strconv"
	"syscall"
)

// 分析进程退出状态
func (session *JudgeSession) analysisExitStatus(rst *TestCaseResult, pinfo *ProcessInfo, specialJudge bool) error {
	ru := pinfo.Rusage
	status := pinfo.Status

	rst.TimeUsed = int(ru.Utime.Sec * 1000 + int64(ru.Utime.Usec) / 1000 + ru.Stime.Sec * 1000 + int64(ru.Stime.Usec) / 1000)
	rst.MemoryUsed = int(ru.Minflt * int64(syscall.Getpagesize() / 1024 ))

	// 特判
	if specialJudge {
		if status.Signaled() {
			sig := status.Signal()
			rst.ReSignum = int(sig)
			if session.SpecialJudge.Mode != SpecialJudgeModeInteractive {
				// 检查判题程序是否超时
				if sig == syscall.SIGXCPU || sig == syscall.SIGALRM {
					rst.JudgeResult = JudgeFlagSpecialJudgeTimeout
					rst.ReInfo = fmt.Sprintf("special judger time limit exceed, unix singal: %d", sig)
				} else {
					rst.JudgeResult = JudgeFlagSpecialJudgeError
					rst.ReInfo = fmt.Sprintf("special judger caused an error, unix singal: %d", sig)
				}
			} else {
				// 交互特判时，如果选手程序让判题程序挂了，视作RE
				rst.JudgeResult = JudgeFlagRE
				rst.ReInfo = fmt.Sprintf("special judger caused an error, unix singal: %d", sig)
			}
		} else if status.Exited() {
			// 如果特判程序正常退出
			exitcode := status.ExitStatus()
			rst.SPJExitCode = exitcode
			// 判断退出代码是否正确
			if exitcode == JudgeFlagAC || exitcode == JudgeFlagPE ||
				exitcode == JudgeFlagWA || exitcode == JudgeFlagOLE ||
				exitcode == JudgeFlagSpecialJudgeRequireChecker {
				rst.JudgeResult = exitcode
			} else {
				rst.JudgeResult = JudgeFlagSpecialJudgeError
				rst.ReInfo = fmt.Sprintf("special judger return with a wrong exitcode: %d", exitcode)
			}
		}
	} else {
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
					if r, e := SignumMap[rst.ReSignum]; e {
						rst.ReInfo = fmt.Sprintf("%s: %s", r[0], r[1])
					}
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
				if r, e := SignumMap[rst.ReSignum]; e {
					rst.ReInfo = fmt.Sprintf("%s: %s", r[0], r[1])
				}
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
	}
	return nil
}


// 基于JudgeOptions进行评测调度
func (session *JudgeSession) judgeOnce(judgeResult *TestCaseResult) error {
	switch session.SpecialJudge.Mode {
	case SpecialJudgeModeDisabled:
		pinfo, err := session.runNormalJudge(judgeResult)
		if err != nil {
			judgeResult.JudgeResult = JudgeFlagSE
			judgeResult.SeInfo = err.Error()
			return err
		}

		err = session.analysisExitStatus(judgeResult, pinfo, false)
		if err != nil {
			return err
		}
	case SpecialJudgeModeChecker, SpecialJudgeModeInteractive:
		tinfo, jinfo, err := session.runSpecialJudge(judgeResult)
		if err != nil {
			judgeResult.JudgeResult = JudgeFlagSE
			judgeResult.SeInfo = err.Error()
			return err
		}
		// 分析目标程序的状态
		err = session.analysisExitStatus(judgeResult, tinfo, false)
		if err != nil {
			return err
		}
		// 如果目标程序正常退出，则去分析判题机
		// 否则，优先输出目标程序的运行情况
		if judgeResult.JudgeResult == 0 {
			err = session.analysisExitStatus(judgeResult, jinfo, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 对一组测试数据运行一次评测
func (session *JudgeSession)runOneCase(tc TestCase, Id string) *TestCaseResult {

	tcResult := TestCaseResult{}
	tcResult.Id = Id
	// 创建相关的文件路径
	tcResult.TestCaseIn = tc.TestCaseIn
	tcResult.TestCaseOut = tc.TestCaseOut
	tcResult.ProgramOut = path.Join(session.SessionDir, Id + "_program.out")
	tcResult.ProgramError = path.Join(session.SessionDir,  Id + "_program.err")
	tcResult.ProgramLog = path.Join(session.SessionDir,  Id + "_program.log")
	tcResult.JudgerOut = path.Join(session.SessionDir,  Id + "_judger.out")
	tcResult.JudgerError = path.Join(session.SessionDir,  Id + "_judger.err")
	tcResult.JudgerLog = path.Join(session.SessionDir,  Id + "_judger.log")
	tcResult.JudgerReport = path.Join(session.SessionDir,  Id + "_judger.report")

	// 运行judge程序
	err := session.judgeOnce(&tcResult)
	if err != nil {
		tcResult.JudgeResult = JudgeFlagSE
		tcResult.SeInfo = err.Error()
	}

	return  &tcResult
}

// 执行评测
func (session *JudgeSession)RunJudge() (JudgeResult, error) {
	judgeResult := JudgeResult{}

	err := session.compileTargetProgram(&judgeResult)
	if err != nil {
		return judgeResult, err
	}

	for i := 0; i < len(session.TestCases); i++ {
		if session.TestCases[i].Id == "" {
			session.TestCases[i].Id = strconv.Itoa(i)
		}
		id := session.TestCases[i].Id
		tcResult := session.runOneCase(session.TestCases[i], id)
		judgeResult.TestCases = append(judgeResult.TestCases, *tcResult)
		if tcResult.JudgeResult == JudgeFlagSE {
			judgeResult.JudgeResult = JudgeFlagSE
			judgeResult.SeInfo = fmt.Sprintf("testcase %s caused a problem", id)
			// SE的情况下直接终止执行
			return judgeResult, err
		}
	}

	return judgeResult, nil
}
