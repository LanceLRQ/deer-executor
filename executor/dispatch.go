package executor

import (
	"fmt"
	"log"
	"os"
	"strconv"
)


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
		session.saveExitRusage(judgeResult, pinfo, false)
		// 分析目标程序的状态
		session.analysisExitStatus(judgeResult, pinfo, false)
		// 只有AC的时候才进行文本比较！
		if judgeResult.JudgeResult == JudgeFlagAC {
			// 进行文本比较
			err = session.DiffText(judgeResult)
			if err != nil {
				judgeResult.JudgeResult = JudgeFlagSE
				judgeResult.SeInfo = err.Error()
				return err
			}
		}

	case SpecialJudgeModeChecker, SpecialJudgeModeInteractive:
		tinfo, jinfo, err := session.runSpecialJudge(judgeResult)
		if err != nil {
			judgeResult.JudgeResult = JudgeFlagSE
			judgeResult.SeInfo = err.Error()
			return err
		}
		session.saveExitRusage(judgeResult, tinfo, false)
		session.saveExitRusage(judgeResult, jinfo, true)
		// 分析判题程序的状态
		session.analysisExitStatus(judgeResult, jinfo, true)
		// 如果判题程序正常退出，则再去分析目标程序
		if judgeResult.JudgeResult == 0 {
			session.analysisExitStatus(judgeResult, tinfo, false)
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
	tcResult.ProgramOut = Id + "_program.out"
	tcResult.ProgramError = Id + "_program.err"
	tcResult.JudgerOut = Id + "_judger.out"
	tcResult.JudgerError = Id + "_judger.err"
	tcResult.JudgerReport = Id + "_judger.report"

	// 运行judge程序
	err := session.judgeOnce(&tcResult)
	if err != nil {
		tcResult.JudgeResult = JudgeFlagSE
		tcResult.SeInfo = err.Error()
	}

	return  &tcResult
}

// 执行评测
func (session *JudgeSession)RunJudge() JudgeResult {
	judgeResult := JudgeResult{}
	judgeResult.SessionId = session.SessionId

	err := session.compileTargetProgram(&judgeResult)
	if err != nil {
		return judgeResult
	}

	if session.SpecialJudge.Mode > 0 {
		// 如果需要特殊评测，则编译相关代码
		err := session.compileJudgerProgram(&judgeResult)
		if err != nil {
			return judgeResult
		}
	}
	tl, ml, rtl, fsl, mlf := getLimitation(session)
	mlfText := ""
	if mlf > 0 {
		mlfText = fmt.Sprintf(" (with %d KB for VM)", mlf)
	}
	log.Printf(
		"Time limit: %d ms, Memory limit: %d KB%s, Real-time limit: %d ms, File size limit: %d KB\n",
		tl, ml, mlfText, rtl, fsl/1024,
	)

	exitCodes := make([]int, 0, 1)
	for i := 0; i < len(session.TestCases); i++ {
		if session.TestCases[i].Id == "" {
			session.TestCases[i].Id = strconv.Itoa(i)
		}
		id := session.TestCases[i].Id

		tcResult := session.runOneCase(session.TestCases[i], id)

		isFault := session.isDisastrousFault(&judgeResult, tcResult)
		judgeResult.TestCases = append(judgeResult.TestCases, *tcResult)
		judgeResult.MemoryUsed = Max32(tcResult.MemoryUsed, judgeResult.MemoryUsed)
		judgeResult.TimeUsed = Max32(tcResult.TimeUsed, judgeResult.TimeUsed)
		// 这里使用动态增加的方式是为了保证len(exitCodes)<=len(testCases)
		// 方便计算最终结果的时候判定测试数据是否全部跑完
		exitCodes = append(exitCodes, tcResult.JudgeResult)

		// 如果发生灾难性错误，直接退出
		if isFault {
			break
		}

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
	session.generateFinallyResult(&judgeResult, exitCodes)
	return judgeResult
}

// 清理案发现场
func (session *JudgeSession)Clean() {
	_ = os.RemoveAll(session.SessionDir)
}