package executor

import (
    "context"
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "os"
    "path"
    "strconv"
    "time"
)

// 基于JudgeOptions进行评测调度
func (session *JudgeSession) judgeOnce(judgeResult *commonStructs.TestCaseResult) {
    switch session.JudgeConfig.SpecialJudge.Mode {
    case constants.SpecialJudgeModeDisabled:
        pinfo, err := session.runNormalJudge(judgeResult)
        if err != nil {
            judgeResult.JudgeResult = constants.JudgeFlagSE
            judgeResult.SeInfo = err.Error()
            return
        }
        session.saveExitRusage(judgeResult, pinfo, false)
        // 分析目标程序的状态
        session.analysisExitStatus(judgeResult, pinfo, false)
        // 只有AC的时候才进行文本比较！
        if judgeResult.JudgeResult == constants.JudgeFlagAC {
            // 进行文本比较
            err = session.DiffText(judgeResult)
            if err != nil {
                judgeResult.JudgeResult = constants.JudgeFlagSE
                judgeResult.SeInfo = err.Error()
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
                // 进行文本比较
                err = session.DiffText(judgeResult)
                if err != nil {
                    judgeResult.JudgeResult = constants.JudgeFlagSE
                    judgeResult.SeInfo = err.Error()
                    return
                }
            }
        }
    }
    return
}


func checkTestCaseInputOutput(tcase commonStructs.TestCase, configDir string, input, output bool) error {
    // 检查Input、Output是否存在
    _, err := os.Stat(path.Join(configDir, tcase.Input))
    if os.IsNotExist(err) {
        return fmt.Errorf("test case (%s) input file (%s) not exists", tcase.Handle, tcase.Input)
    }
    _, err = os.Stat(path.Join(configDir, tcase.Output))
    if os.IsNotExist(err) {
        return fmt.Errorf("test case (%s) output file (%s) not exists", tcase.Handle, tcase.Output)
    }
    return nil
}


// 对一组测试数据运行一次评测
func (session *JudgeSession) runOneCase(config *commonStructs.JudgeConfiguration, tc commonStructs.TestCase, Id string) *commonStructs.TestCaseResult {

     var err error

    tcResult := commonStructs.TestCaseResult{}
    tcResult.Handle = Id
    // 创建相关的文件路径
    tcResult.Input = tc.Input
    tcResult.Output = tc.Output
    tcResult.ProgramOut = Id + "_program.out"
    tcResult.ProgramError = Id + "_program.err"
    tcResult.JudgerOut = Id + "_judger.out"
    tcResult.JudgerError = Id + "_judger.err"
    tcResult.JudgerReport = Id + "_judger.report"


    if config.SpecialJudge.Mode > 0 && tc.UseGenerator {
        // 当普通特殊评测启用的时候
        // 判断是否有generator，如果有先运行它并获取输入数据
        ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
        _, err = utils.CallGenerator(ctx, &tc, config.ConfigDir)
        if err != nil {
            tcResult.JudgeResult = constants.JudgeFlagSE
            tcResult.SeInfo = err.Error()
            return &tcResult
        }
        // TODO 写入到Input
    }

    if config.SpecialJudge.Mode == 0 {
        // 非特判模式检查Input、Output是否存在
        err =checkTestCaseInputOutput(tc, config.ConfigDir, true, true)
    } else {
        // 特判模式只需要检查Input是否存在
        err = checkTestCaseInputOutput(tc, config.ConfigDir, true, false)
    }
    if err != nil {
        tcResult.JudgeResult = constants.JudgeFlagSE
        tcResult.SeInfo = err.Error()
        return &tcResult
    }

    // 运行judge程序
    session.judgeOnce(&tcResult)

    return &tcResult
}

// 执行评测
func (session *JudgeSession) RunJudge() commonStructs.JudgeResult {
    judgeResult := commonStructs.JudgeResult{}
    judgeResult.SessionId = session.SessionId

    err := session.compileTargetProgram(&judgeResult)
    if err != nil {
        return judgeResult
    }

    if session.JudgeConfig.SpecialJudge.Mode > 0 {
        // 如果需要特殊评测，则编译相关代码
        err := session.compileJudgerProgram(&judgeResult)
        if err != nil {
            return judgeResult
        }
    }
    //tl, ml, rtl, fsl, mlf := getLimitation(session)
    //mlfText := ""
    //if mlf > 0 {
    //	mlfText = fmt.Sprintf(" (with %d KB for VM)", mlf)
    //}
    //log.Printf(
    //	"Time limit: %d ms, Memory limit: %d KB%s, Real-time limit: %d ms, File size limit: %d KB\n",
    //	tl, ml, mlfText, rtl, fsl/1024,
    //)

    exitCodes := make([]int, 0, 1)
    for i := 0; i < len(session.JudgeConfig.TestCases); i++ {
        if session.JudgeConfig.TestCases[i].Handle == "" {
            session.JudgeConfig.TestCases[i].Handle = strconv.Itoa(i)
        }
        id := session.JudgeConfig.TestCases[i].Handle

        tcResult := session.runOneCase(&session.JudgeConfig, session.JudgeConfig.TestCases[i], id)

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
    return judgeResult
}

// 清理案发现场
func (session *JudgeSession) Clean() {
    _ = os.RemoveAll(session.SessionDir)
}
