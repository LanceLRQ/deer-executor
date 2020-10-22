package executor

import (
	"fmt"
	"io/ioutil"
	"syscall"
)

// 分析进程退出状态
func (session *JudgeSession) analysisExitStatus(rst *TestCaseResult, pinfo *ProcessInfo, judger bool) error {
	ru := pinfo.Rusage
	status := pinfo.Status

	tu := int(ru.Utime.Sec*1000 + int64(ru.Utime.Usec)/1000 + ru.Stime.Sec*1000 + int64(ru.Stime.Usec)/1000)
	mu := int(ru.Minflt * int64(syscall.Getpagesize()/1024))

	// 特判
	if judger {
		rst.SPJTimeUsed = tu
		rst.SPJMemoryUsed = mu
		if status.Signaled() {
			sig := status.Signal()
			rst.SPJReSignum = int(sig)
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
		rst.TimeUsed = tu
		rst.MemoryUsed = mu
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
					if r, e := SignalNumberMap[rst.ReSignum]; e {
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
				if r, e := SignalNumberMap[rst.ReSignum]; e {
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

// 判定是否是灾难性结果
func (session *JudgeSession) isDisastrousFault(judgeResult *JudgeResult, tcResult *TestCaseResult) bool {
	if tcResult.JudgeResult == JudgeFlagSE {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = fmt.Sprintf("testcase %s caused a problem", tcResult.Id)
		return true
	}

	// 如果是实时运行的语言
	if session.compiler.IsRealTime() {
		outfile, e := ioutil.ReadFile(tcResult.ProgramError)
		if e == nil {
			if len(outfile) > 0 {
				remsg := string(outfile)

				if session.compiler.IsCompileError(remsg) {
					tcResult.JudgeResult = JudgeFlagCE
					tcResult.CeInfo = remsg
					judgeResult.JudgeResult = JudgeFlagCE
					judgeResult.CeInfo = remsg
				} else {
					tcResult.JudgeResult = JudgeFlagRE
					tcResult.SeInfo = fmt.Sprintf("%s\n%s\n", tcResult.SeInfo, remsg)
					judgeResult.JudgeResult = JudgeFlagRE
					judgeResult.SeInfo = tcResult.SeInfo
				}
				return true
			}
		}
	}
	return false
}

// 计算判题结果
func (session *JudgeSession) generateFinallyResult(result *JudgeResult, exitcodes []int) {
	var (
		ac, pe, wa int = 0, 0, 0
	)
	for _, exitcode := range exitcodes {
		// 如果，不是AC、PE、WA
		if exitcode != JudgeFlagWA && exitcode != JudgeFlagPE && exitcode != JudgeFlagAC {
			//直接应用结果
			result.JudgeResult = exitcode
			return
		}
		if exitcode == JudgeFlagWA { wa++ }
		if exitcode == JudgeFlagPE { pe++ }
		if exitcode == JudgeFlagAC { ac++ }
	}
	// 在严格判题模式下，由于第一组数据不是AC\PE就会直接报错，因此要判定测试数据是否全部跑完。
	if len(exitcodes) != len(session.TestCases) {
		// 如果测试数据未全部跑完
		result.JudgeResult = JudgeFlagWA
	} else {
		// 如果测试数据未全部跑了
		if wa > 0 {
			// 如果存在WA，报WA
			result.JudgeResult = JudgeFlagWA
		} else if pe > 0 {	// 如果PE > 0
			if !session.StrictMode {
				// 非严格模式，报AC
				result.JudgeResult = JudgeFlagAC
			} else {
				// 严格模式下报PE
				result.JudgeResult = JudgeFlagPE
			}
		} else {
			result.JudgeResult = JudgeFlagAC
		}
	}
}