package packmgr

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/executor"
    "github.com/urfave/cli/v2"
    "io/ioutil"
    "log"
    "os"
    "path"
)

func runCheckerCase(session *executor.JudgeSession, caseIndex int) error {
    iCase := &session.JudgeConfig.SpecialJudge.CheckerCases[caseIndex]


    // 写入数据
    err := ioutil.WriteFile(path.Join(session.SessionDir, string(caseIndex) + "_problem.in"), []byte(iCase.Input), 0666)
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(path.Join(session.SessionDir, string(caseIndex) + "_problem.out"), []byte(iCase.Output), 0666)
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(path.Join(session.SessionDir, string(caseIndex) + "_problem.answer"), []byte(iCase.Answer), 0666)
    if err != nil {
        return err
    }

    // TODO 运行checker程序

    return nil
}


func isCheckerExists(config *structs.JudgeConfiguration) error {
    cPath, err := utils.GetCompiledBinaryFileAbsPath("checker", config.SpecialJudge.Name, config.ConfigDir)
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
            err := runCheckerCase(session, key)
            if err != nil { return err }
        }
    } else {
        log.Printf("[generator] run case #%d", caseIndex)
        err := runCheckerCase(session, caseIndex)
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

    if session.JudgeConfig.SpecialJudge.Mode != 1 {
        return fmt.Errorf("[checker] only support checker mode")
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