package executor

import (
	"fmt"
	"io/ioutil"
	"os"
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
		// 分析目标程序的状态
		err = session.analysisExitStatus(judgeResult, pinfo, false)
		if err != nil {
			return err
		}
		// 进行文本比较
		err = session.DiffText(judgeResult)
		if err != nil {
			judgeResult.JudgeResult = JudgeFlagSE
			judgeResult.SeInfo = err.Error()
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
		// 普通checker的时候支持按判题机的意愿进行文本比较
		if session.SpecialJudge.Mode == SpecialJudgeModeChecker {
			if judgeResult.JudgeResult == JudgeFlagSpecialJudgeRequireChecker {
				// 进行文本比较
				err = session.DiffText(judgeResult)
				if err != nil {
					judgeResult.JudgeResult = JudgeFlagSE
					judgeResult.SeInfo = err.Error()
					return err
				}
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
			if session.StrictMode {
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

// 执行评测
func (session *JudgeSession)RunJudge() (JudgeResult, error) {
	judgeResult := JudgeResult{}

	err := session.compileTargetProgram(&judgeResult)
	if err != nil {
		return judgeResult, err
	}
	exitcodes := make([]int, 0, 1)
	for i := 0; i < len(session.TestCases); i++ {
		if session.TestCases[i].Id == "" {
			session.TestCases[i].Id = strconv.Itoa(i)
		}
		id := session.TestCases[i].Id

		tcResult := session.runOneCase(session.TestCases[i], id)

		isFalut := session.isDisastrousFault(&judgeResult, tcResult)
		judgeResult.TestCases = append(judgeResult.TestCases, *tcResult)
		judgeResult.MemoryUsed = Max32(tcResult.MemoryUsed, judgeResult.MemoryUsed)
		judgeResult.TimeUsed = Max32(tcResult.TimeUsed, judgeResult.TimeUsed)
		// 如果发生灾难性错误，直接退出
		if isFalut {
			break
		}
		// 这里使用动态增加的方式是为了保证len(exitcodes)<=len(testCases)
		// 方便计算最终结果的时候判定测试数据是否全部跑完
		exitcodes = append(exitcodes, tcResult.JudgeResult)

		//判定是否继续判题
		keep := false
		if tcResult.JudgeResult == JudgeFlagAC || tcResult.JudgeResult == JudgeFlagPE {
			keep = true
		} else if !session.StrictMode && tcResult.JudgeResult == JudgeFlagWA {
			keep = true
		}
		if !keep {
			break
		}
	}
	// 计算最终结果
	session.generateFinallyResult(&judgeResult, exitcodes)

	return judgeResult, nil
}

func (session *JudgeSession)Clean() {
	_ = os.RemoveAll(session.SessionDir)
}