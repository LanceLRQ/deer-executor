package packmgr

import (
    "context"
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/v2/executor"
    "github.com/pkg/errors"
    "github.com/urfave/cli/v2"
    "io/ioutil"
    "log"
    "os"
    "path"
    "strings"
    "time"
)

func runCheckerCase(session *executor.JudgeSession, caseIndex int) error {
    iCase := &session.JudgeConfig.SpecialJudge.CheckerCases[caseIndex]
    tInput := path.Join(session.SessionDir, fmt.Sprintf("%d_problem.in", caseIndex))
    tOutput := path.Join(session.SessionDir, fmt.Sprintf("%d_problem.out", caseIndex))
    tAnswer := path.Join(session.SessionDir, fmt.Sprintf("%d_problem.ans", caseIndex))

    // 写入数据
    err := ioutil.WriteFile(tInput, []byte(iCase.Input), 0666)
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(tOutput, []byte(iCase.Output), 0666)
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(tAnswer, []byte(iCase.Answer), 0666)
    if err != nil {
        return err
    }

    config := session.JudgeConfig

    cPath, err := utils.GetCompiledBinaryFileAbsPath("checker", config.SpecialJudge.Name, session.ConfigDir)
    if err != nil { return err }
    // 运行checker程序
    ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
    // <input-file> <output-file> <answer-file> [<report-file>]
    ret, err := utils.RunUnixShell(&structs.ShellOptions{
        Context:   ctx,
        Name:      cPath,
        Args:      []string {
            tInput,
            tOutput,
            tAnswer,
        },
        StdWriter: nil,
        OnStart:   nil,
    })
    if err != nil { return err }

    judgeResult := -1
    for _, item := range constants.TestlibExitMsgMapping {
        if strings.HasPrefix(ret.Stderr, item.ErrName) {
            judgeResult = item.JudgeResult
        }
    }
    if judgeResult == -1 {
        judgeResult = constants.JudgeFlagSpecialJudgeError
    }
    iCase.CheckerVerdict = judgeResult
    iCase.CheckerComment = ret.Stderr
    iCase.Verdict = iCase.CheckerVerdict == iCase.ExpectedVerdict
    return nil
}


// 检查是否存在checker
func isCheckerExists(config *structs.JudgeConfiguration) error {
    cPath, err := utils.GetCompiledBinaryFileAbsPath("checker", config.SpecialJudge.Name, config.ConfigDir)
    if err != nil { return err }
    _, err = os.Stat(cPath)
    if os.IsNotExist(err) {
        return errors.Errorf("[checker] execuable checker file not exists")
    }
    return nil
}


// 遍历checker cases
func runCheckerCases (session *executor.JudgeSession, caseIndex int) error {
    if caseIndex < 0 {
        for key, _ := range session.JudgeConfig.SpecialJudge.CheckerCases {
            log.Printf("[checker] run case #%d", key)
            err := runCheckerCase(session, key)
            if err != nil { return err }
        }
    } else {
        log.Printf("[checker] run case #%d", caseIndex)
        err := runCheckerCase(session, caseIndex)
        if err != nil { return err }
    }
    return nil
}

// 运行特殊评测的checker (APP入口)
func RunCheckerCases(c *cli.Context) error {
    configFile := c.Args().Get(0)
    _, err := os.Stat(configFile)
    if err != nil && os.IsNotExist(err) {
        return errors.Errorf("[checker] problem config file (%s) not found", configFile)
    }
    session, err := executor.NewSession(configFile)
    if err != nil { return err }

    if session.JudgeConfig.SpecialJudge.Mode != 1 {
        return errors.Errorf("[checker] only support checker mode")
    }

    // checker exists?
    if err = isCheckerExists(&session.JudgeConfig); err != nil {
        return err
    }

    silence := c.Bool("silence")
    answerCaseIndex := c.Uint("answer")
    caseIndex := c.Int("case")

    // Init work dir and compile answer
    err = initWork(session, answerCaseIndex)
    if err != nil {
        return err
    }
    defer session.Clean()

    err = runCheckerCases(session, caseIndex)
    if err != nil {
        return err
    }
    return session.SaveConfiguration(!silence)
}