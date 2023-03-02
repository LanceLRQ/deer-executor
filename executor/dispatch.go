//go:build linux || darwin
// +build linux darwin

package executor

import (
	"github.com/LanceLRQ/deer-executor/v2/common/constants"
	commonStructs "github.com/LanceLRQ/deer-executor/v2/common/structs"
	"github.com/pkg/errors"
	"os"
	"path"
	"strconv"
)

// JudgeOnce 基于JudgeOptions进行评测调度
func (session *JudgeSession) JudgeOnce(judgeResult *commonStructs.TestCaseResult) {
	switch session.JudgeConfig.SpecialJudge.Mode {
	case constants.SpecialJudgeModeDisabled:
		pinfo, err := session.runNormalJudge(judgeResult)
		if err != nil {
			judgeResult.JudgeResult = constants.JudgeFlagSE
			judgeResult.SeInfo = err.Error()
			session.Logger.Error(err.Error())
			return
		}
		session.saveExitRusage(judgeResult, pinfo, false)
		// 分析目标程序的状态
		session.analysisExitStatus(judgeResult, pinfo, false)
		// 只有AC的时候才进行文本比较！
		if judgeResult.JudgeResult == constants.JudgeFlagAC {
			session.Logger.Infof("Run text checker.")
			// 进行文本比较
			err = session.DiffText(judgeResult)
			if err != nil {
				judgeResult.JudgeResult = constants.JudgeFlagSE
				judgeResult.SeInfo = err.Error()
				session.Logger.Error(err.Error())
				return
			}
		}

	case constants.SpecialJudgeModeChecker, constants.SpecialJudgeModeInteractive:
		tinfo, jinfo, err := session.runSpecialJudge(judgeResult)
		if err != nil {
			judgeResult.JudgeResult = constants.JudgeFlagSE
			judgeResult.SeInfo = err.Error()
			return
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
		if session.JudgeConfig.SpecialJudge.Mode == constants.SpecialJudgeModeChecker {
			if judgeResult.JudgeResult == constants.JudgeFlagSpecialJudgeRequireChecker {
				session.Logger.Infof("Run text checker.")
				// 进行文本比较
				err = session.DiffText(judgeResult)
				if err != nil {
					judgeResult.JudgeResult = constants.JudgeFlagSE
					judgeResult.SeInfo = err.Error()
					session.Logger.Error(err.Error())
					return
				}
			}
		}
	}
	return
}

// 检查Input、Output是否存在
func checkTestCaseInputOutput(tcase commonStructs.TestCase, configDir string) error {
	_, err := os.Stat(path.Join(configDir, tcase.Input))
	if os.IsNotExist(err) {
		return errors.Errorf("test case (%s) input file (%s) not exists", tcase.Handle, tcase.Input)
	}
	_, err = os.Stat(path.Join(configDir, tcase.Output))
	if os.IsNotExist(err) {
		return errors.Errorf("test case (%s) output file (%s) not exists", tcase.Handle, tcase.Output)
	}
	return nil
}

// 对一组测试数据运行一次评测
func (session *JudgeSession) runOneCase(config *commonStructs.JudgeConfiguration, tc commonStructs.TestCase, id string) *commonStructs.TestCaseResult {
	session.Logger.Infof("Run test case: %s", id)

	var err error

	tcResult := commonStructs.TestCaseResult{}
	tcResult.Handle = id
	// 创建相关的文件路径
	tcResult.Input = tc.Input
	tcResult.Output = tc.Output
	tcResult.ProgramOut = id + "_program.out"
	tcResult.ProgramError = id + "_program.err"
	tcResult.CheckerOut = id + "_checker.out"
	tcResult.CheckerError = id + "_checker.err"
	tcResult.CheckerReport = id + "_checker.report"

	// 检查测试数据的输入输出文件是否存在
	err = checkTestCaseInputOutput(tc, config.ConfigDir)
	if err != nil {
		tcResult.JudgeResult = constants.JudgeFlagSE
		tcResult.SeInfo = err.Error()
		session.Logger.Error(err.Error())
		return &tcResult
	}

	// 运行judge程序
	session.JudgeOnce(&tcResult)

	return &tcResult
}

// RunJudge 执行评测
func (session *JudgeSession) RunJudge() commonStructs.JudgeResult {
	session.Logger.Info("Start Judgement")

	// make judge result
	judgeResult := commonStructs.JudgeResult{}
	judgeResult.SessionID = session.SessionID

	// compile code
	err := session.compileTargetProgram(&judgeResult)
	if err != nil {
		judgeResult.JudgeLogs = session.Logger.GetLogs()
		return judgeResult
	}

	if session.JudgeConfig.SpecialJudge.Mode > 0 {
		// 如果需要特殊评测，则编译相关代码
		err := session.compileJudgerProgram(&judgeResult)
		if err != nil {
			judgeResult.JudgeResult = constants.JudgeFlagSE
			judgeResult.SeInfo = err.Error()
			judgeResult.JudgeLogs = session.Logger.GetLogs()
			return judgeResult
		}
	}

	// 资源限制信息更新
	updateLimitation(session)

	session.Logger.Info("Ready for judgement")
	// Init exit code
	exitCodes := make([]int, 0, 1)
	for i := 0; i < len(session.JudgeConfig.TestCases); i++ {
		if session.JudgeConfig.TestCases[i].Handle == "" {
			session.JudgeConfig.TestCases[i].Handle = strconv.Itoa(i)
		}
		id := session.JudgeConfig.TestCases[i].Handle

		tcResult := session.runOneCase(&session.JudgeConfig, session.JudgeConfig.TestCases[i], id)

		flagName, ok := constants.FlagMeansMap[tcResult.JudgeResult]
		if !ok {
			flagName = "Unknown"
		}
		session.Logger.Infof("This case's result is " + flagName)

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
		if tcResult.JudgeResult == constants.JudgeFlagAC || tcResult.JudgeResult == constants.JudgeFlagPE {
			keep = true
		} else if !session.JudgeConfig.StrictMode && tcResult.JudgeResult == constants.JudgeFlagWA {
			keep = true
		}
		if !keep {
			break
		}
	}
	// 计算最终结果
	session.generateFinallyResult(&judgeResult, exitCodes)

	// Log
	if judgeResult.JudgeResult == constants.JudgeFlagAC {
		session.Logger.Info("Congratulations! Your code have been ACCEPTED.")
	} else {
		flagName, ok := constants.FlagMeansMap[judgeResult.JudgeResult]
		if !ok {
			flagName = "Unknown"
		}
		session.Logger.Warnf("Oops, %s! There is something wrong.", flagName)
	}
	// get judge logs
	judgeResult.JudgeLogs = session.Logger.GetLogs()
	// return
	return judgeResult
}

// 从资源限制的参数列表里按语言获取相关信息，并作为当前的资源限制参数。
func updateLimitation(session *JudgeSession) {
	langName := session.Compiler.GetName()
	memoryLimitExtend := 0
	jitMem, ok := constants.MemorySizeForJIT[langName]
	if ok {
		memoryLimitExtend = jitMem
	}
	limitation, ok := session.JudgeConfig.Limitation[langName]
	if ok {
		session.JudgeConfig.TimeLimit = limitation.TimeLimit
		session.JudgeConfig.MemoryLimit = limitation.MemoryLimit + memoryLimitExtend
		session.JudgeConfig.RealTimeLimit = limitation.RealTimeLimit
		session.JudgeConfig.FileSizeLimit = limitation.FileSizeLimit
		return
	}
	session.JudgeConfig.MemoryLimit = session.JudgeConfig.MemoryLimit + memoryLimitExtend
}
