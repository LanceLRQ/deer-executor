package executor

import (
	"os"
	"path"
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
		// 分析目标程序的状态
		err = session.analysisExitStatus(judgeResult, pinfo, false)
		if err != nil {
			return err
		}
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
	tcResult.JudgerOut = path.Join(session.SessionDir,  Id + "_judger.out")
	tcResult.JudgerError = path.Join(session.SessionDir,  Id + "_judger.err")
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