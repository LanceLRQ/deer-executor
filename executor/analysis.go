//go:build linux || darwin
// +build linux darwin

package executor

import (
	"fmt"
	constants2 "github.com/LanceLRQ/deer-executor/v3/executor/constants"
	commonStructs "github.com/LanceLRQ/deer-executor/v3/executor/structs"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"io/ioutil"
	"path"
	"strconv"
	"syscall"
)

// 分析进程资源占用
func (session *JudgeSession) saveExitRusage(rst *commonStructs.TestCaseResult, pinfo *ProcessInfo, judger bool) {
	ru := pinfo.Rusage
	status := pinfo.Status

	tu := int(ru.Utime.Sec*1000 + int64(ru.Utime.Usec)/1000 + ru.Stime.Sec*1000 + int64(ru.Stime.Usec)/1000)
	//maxrss := ru.Maxrss
	//if runtime.GOOS == "darwin" {
	//    maxrss = maxrss / 1024
	//}
	mu := int(ru.Minflt * int64(syscall.Getpagesize()/1024))

	// 特判
	if judger {
		rst.SPJTimeUsed = tu
		rst.SPJMemoryUsed = mu
		rst.SPJReSignum = int(status.Signal())
		session.Logger.Infof(
			"checker exit with code: %d, signum: %d, Time used: %d, Mem used: %d.",
			pinfo.Status.ExitStatus(),
			rst.SPJReSignum,
			rst.SPJTimeUsed,
			rst.SPJMemoryUsed,
		)
	} else {
		rst.TimeUsed = tu
		rst.MemoryUsed = mu
		rst.ReSignum = int(status.Signal())
		session.Logger.Infof(
			"program exit with code: %d, signum: %d, Time used: %d, Mem used: %d.",
			pinfo.Status.ExitStatus(),
			rst.ReSignum,
			rst.TimeUsed,
			rst.MemoryUsed,
		)
	}
}

// 分析进程退出状态
func (session *JudgeSession) analysisExitStatus(rst *commonStructs.TestCaseResult, pinfo *ProcessInfo, judger bool) {
	status := pinfo.Status

	// 特判
	if judger {
		if status.Signaled() {
			sig := status.Signal()
			if session.JudgeConfig.SpecialJudge.Mode != constants2.SpecialJudgeModeInteractive {
				// 检查判题程序是否超时
				if sig == syscall.SIGXCPU || sig == syscall.SIGALRM {
					rst.JudgeResult = constants2.JudgeFlagSpecialJudgeTimeout
					rst.ReInfo = fmt.Sprintf("special judger time limit exceed, unix singal: %d", sig)
				} else {
					rst.JudgeResult = constants2.JudgeFlagSpecialJudgeError
					rst.ReInfo = fmt.Sprintf("special judger caused an error, unix singal: %d", sig)
				}
			} else {
				// 交互特判时，如果选手程序让判题程序挂了，视作RE
				rst.JudgeResult = constants2.JudgeFlagRE
				rst.ReInfo = fmt.Sprintf("special judger caused an error, unix singal: %d", sig)
			}
		} else if status.Exited() {
			// 如果特判程序正常退出
			exitcode := status.ExitStatus()
			rst.SPJExitCode = exitcode
			if session.JudgeConfig.SpecialJudge.UseTestlib {
				// 如果是Testlib的checker，则退出代码要按照他们的规则去判定
				msg, err := ioutil.ReadFile(path.Join(session.SessionDir, rst.CheckerReport))
				if err != nil {
					rst.JudgeResult = constants2.JudgeFlagSpecialJudgeError
					rst.SPJMsg = fmt.Sprintf("read checker report file error: %s", err.Error())
				} else {
					tr := commonStructs.TestlibCheckerResult{}
					ok := utils.XMLStringObject(string(msg), &tr)
					if ok {
						rst.SPJMsg = tr.Description
						if flag, ok := constants2.TestlibOutcomeMapping[tr.Outcome]; ok {
							rst.JudgeResult = flag
						} else {
							rst.JudgeResult = constants2.JudgeFlagSpecialJudgeError
						}
						if tr.Outcome == "partially-correct" {
							rst.PartiallyScore, _ = strconv.Atoi(tr.PcType)
						}
					} else {
						rst.JudgeResult = constants2.JudgeFlagSpecialJudgeError
						rst.SPJMsg = fmt.Sprintf("parsee checker report file error:\n%s", string(msg))
					}
				}
			} else {
				// 判断退出代码是否正确
				if exitcode == constants2.JudgeFlagAC || exitcode == constants2.JudgeFlagPE ||
					exitcode == constants2.JudgeFlagWA || exitcode == constants2.JudgeFlagOLE ||
					exitcode == constants2.JudgeFlagSpecialJudgeRequireChecker {
					rst.JudgeResult = exitcode
				} else {
					rst.JudgeResult = constants2.JudgeFlagSpecialJudgeError
					rst.SPJMsg = fmt.Sprintf("special judger return with a wrong exitcode: %d", exitcode)
				}
			}
		}
	} else {
		// If process stopped with a signal
		if status.Signaled() {
			sig := status.Signal()
			if sig == syscall.SIGSEGV {
				// MLE or RE can also get SIGSEGV signal.
				if rst.MemoryUsed > session.JudgeConfig.MemoryLimit {
					rst.JudgeResult = constants2.JudgeFlagMLE
				} else {
					rst.JudgeResult = constants2.JudgeFlagRE
					if r, e := constants2.SignalNumberMap[rst.ReSignum]; e {
						rst.ReInfo = fmt.Sprintf("%s: %s", r[0], r[1])
					}
				}
			} else if sig == syscall.SIGXFSZ {
				// SIGXFSZ signal means OLE
				rst.JudgeResult = constants2.JudgeFlagOLE
			} else if sig == syscall.SIGALRM || sig == syscall.SIGVTALRM || sig == syscall.SIGXCPU {
				// Normal TLE signal
				rst.JudgeResult = constants2.JudgeFlagTLE
			} else if sig == syscall.SIGKILL {
				// Sometimes MLE might get SIGKILL signal.
				// So if real time used lower than TIME_LIMIT - 100, it might be a TLE error.
				if rst.TimeUsed > (session.JudgeConfig.TimeLimit - 100) {
					rst.JudgeResult = constants2.JudgeFlagTLE
				} else {
					rst.JudgeResult = constants2.JudgeFlagMLE
				}
			} else {
				// Otherwise, called runtime error.
				rst.JudgeResult = constants2.JudgeFlagRE
				if r, e := constants2.SignalNumberMap[rst.ReSignum]; e {
					rst.ReInfo = fmt.Sprintf("%s: %s", r[0], r[1])
				}
			}
		} else {
			// Sometimes setrlimit doesn't work accurately.
			if rst.TimeUsed > session.JudgeConfig.TimeLimit {
				rst.JudgeResult = constants2.JudgeFlagMLE
			} else if rst.MemoryUsed > session.JudgeConfig.MemoryLimit {
				rst.JudgeResult = constants2.JudgeFlagMLE
			} else {
				rst.JudgeResult = constants2.JudgeFlagAC
			}
		}
	}
}

// 判定是否是灾难性结果
func (session *JudgeSession) isDisastrousFault(judgeResult *commonStructs.JudgeResult, tcResult *commonStructs.TestCaseResult) bool {
	if tcResult.JudgeResult == constants2.JudgeFlagSE {
		judgeResult.JudgeResult = constants2.JudgeFlagSE
		judgeResult.SeInfo = fmt.Sprintf("testcase %s caused a problem", tcResult.Handle)
		return true
	}

	// 如果是实时运行的语言
	if session.Compiler.IsRealTime() {
		outfile, e := ioutil.ReadFile(path.Join(session.SessionDir, tcResult.ProgramError))
		if e == nil {
			if len(outfile) > 0 {
				remsg := string(outfile)

				if session.Compiler.IsCompileError(remsg) {
					tcResult.JudgeResult = constants2.JudgeFlagCE
					tcResult.CeInfo = remsg
					judgeResult.JudgeResult = constants2.JudgeFlagCE
					judgeResult.CeInfo = remsg
				} else {
					tcResult.JudgeResult = constants2.JudgeFlagRE
					tcResult.SeInfo = fmt.Sprintf("%s\n%s\n", tcResult.SeInfo, remsg)
					judgeResult.JudgeResult = constants2.JudgeFlagRE
					judgeResult.SeInfo = tcResult.SeInfo
				}
				return true
			}
		}
	}
	return false
}

// 计算判题结果
func (session *JudgeSession) generateFinallyResult(result *commonStructs.JudgeResult, exitcodes []int) {
	var (
		ac, pe, wa int = 0, 0, 0
	)
	for _, exitcode := range exitcodes {
		// 如果，不是AC、PE、WA
		if exitcode != constants2.JudgeFlagWA && exitcode != constants2.JudgeFlagPE && exitcode != constants2.JudgeFlagAC {
			//直接应用结果
			result.JudgeResult = exitcode
			return
		}
		if exitcode == constants2.JudgeFlagWA {
			wa++
		}
		if exitcode == constants2.JudgeFlagPE {
			pe++
		}
		if exitcode == constants2.JudgeFlagAC {
			ac++
		}
	}
	// 在严格判题模式下，由于第一组数据不是AC\PE就会直接报错，因此要判定测试数据是否全部跑完。
	if len(exitcodes) != len(session.JudgeConfig.TestCases) {
		// 如果测试数据未全部跑完
		result.JudgeResult = constants2.JudgeFlagWA
	} else {
		// 如果测试数据未全部跑了
		if wa > 0 {
			// 如果存在WA，报WA
			result.JudgeResult = constants2.JudgeFlagWA
		} else if pe > 0 { // 如果PE > 0
			if !session.JudgeConfig.StrictMode {
				// 非严格模式，报AC
				result.JudgeResult = constants2.JudgeFlagAC
			} else {
				// 严格模式下报PE
				result.JudgeResult = constants2.JudgeFlagPE
			}
		} else {
			result.JudgeResult = constants2.JudgeFlagAC
		}
	}
}
