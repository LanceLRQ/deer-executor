package packmgr

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/executor"
    "github.com/urfave/cli/v2"
    "log"
    "os"
)

func runSpecialJudge(session *executor.JudgeSession, caseIndex int) error {

    return nil
}


func isCheckerExists(config *structs.JudgeConfiguration) error {
    cType := "checker"
    if config.SpecialJudge.Mode == 2 {
        cType = "interactor"
    }
    cPath, err := utils.GetCompiledBinaryFileAbsPath(cType, config.SpecialJudge.Name, config.ConfigDir)
    if err != nil { return err }
    _, err = os.Stat(cPath)
    if os.IsNotExist(err) {
        return fmt.Errorf("[checker] execuable checker file not exists")
    }
    return nil
}


// 遍历checker cases
func runSpecialJudgeChecker (session *executor.JudgeSession, caseIndex int) error {
    if caseIndex < 0 {
        for key, _ := range session.JudgeConfig.SpecialJudge.CheckerCases {
            log.Printf("[generator] run case #%d", key)
            err := runSpecialJudge(session, key)
            if err != nil { return err }
        }
    } else {
        log.Printf("[generator] run case #%d", caseIndex)
        err := runSpecialJudge(session, caseIndex)
        if err != nil { return err }
    }
    return nil
}

// 运行特殊评测的checker (APP入口)
func RunSpecialJudgeChecker(c *cli.Context) error {
    configFile := c.Args().Get(0)
    _, err := os.Stat(configFile)
    if err != nil && os.IsNotExist(err) {
        return fmt.Errorf("[checker] problem config file (%s) not found", configFile)
    }
    session, err := executor.NewSession(configFile)
    if err != nil { return err }

    if session.JudgeConfig.SpecialJudge.Mode != 1 && session.JudgeConfig.SpecialJudge.Mode != 2 {
        return fmt.Errorf("[checker] problem not support special judge")
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

    err = runSpecialJudgeChecker(session, caseIndex)
    if err != nil {
        return err
    }
    return session.SaveConfiguration(!silence)
}