package packmgr

import (
    "context"
    "fmt"
    "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/executor"
    "github.com/urfave/cli/v2"
    "io"
    "log"
    "os"
    "time"
)

func runValidator(vBin string, vCase *structs.TestlibValidatorCase) error {
    ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
    rel, err := utils.RunUnixShell(ctx, vBin, nil, func(writer io.Writer) error {
        _, err := writer.Write([]byte(vCase.Input))
        return err
    })
    if err != nil { return err }
    if rel.Success {
        vCase.ValidatorVerdict = true
        vCase.ValidatorComment = ""
    } else {
        log.Printf("[validator] validator error: %s", rel.Stderr)
        vCase.ValidatorVerdict = false
        vCase.ValidatorComment = rel.Stderr
    }
    vCase.Verdict = vCase.ValidatorVerdict == vCase.ExpectedVerdict
    return nil
}

func isValidatorExists (config *structs.JudgeConfiguration) error {
    validator, err := utils.GetCompiledBinaryFileAbsPath("validator", config.TestLib.ValidatorName, config.ConfigDir)
    if err != nil { return err }
    _, err = os.Stat(validator)
    if os.IsNotExist(err) {
        return fmt.Errorf("[validator] execuable validator file not exists")
    }
    return nil
}

// 运行validator case的校验
// caseIndex < 0 表示校验全部
func RunTestlibValidatorCases(config *structs.JudgeConfiguration, caseIndex int) error {
    validator, err := utils.GetCompiledBinaryFileAbsPath("validator", config.TestLib.ValidatorName, config.ConfigDir)
    if err != nil { return err }
    // 执行遍历
    if caseIndex < 0 {
        for key, _ := range config.TestLib.ValidatorCases {
            log.Printf("[validator] run case #%d", key)
            err := runValidator(validator, &config.TestLib.ValidatorCases[key])
            if err != nil { return err }
        }
    } else {
        log.Printf("[validator] run case #%d", caseIndex)
        err := runValidator(validator, &config.TestLib.ValidatorCases[caseIndex])
        if err != nil { return err }
    }
    return nil
}

// 运行Testlib的validator校验
func runTestlibValidators(config *structs.JudgeConfiguration, moduleName string, caseIndex int) error {
    if err := isValidatorExists(config); err != nil {
        return err
    }
    if moduleName == "all" {
        caseIndex = -1
    }
    if moduleName == "all" || moduleName == "validate_cases" {
        if err := RunTestlibValidatorCases(config, caseIndex) ; err != nil {
            return err
        }
    }
    // TODO for test_cases
    //if moduleName == "all" || moduleName == "test_cases" {
    //    if err := RunTestlibValidatorCases(config, -1) ; err != nil {
    //        return err
    //    }
    //}
    return nil
}

// 运行Testlib的validator校验 (APP入口)
func RunTestlibValidators(c *cli.Context) error {
    configFile := c.Args().Get(0)
    _, err := os.Stat(configFile)
    if err != nil && os.IsNotExist(err) {
        return fmt.Errorf("[validator] problem config file (%s) not found", configFile)
    }
    session, err := executor.NewSession(configFile)
    if err != nil { return err }
    mtype := c.String("type")
    mCaseIndex := c.Int("case")

    LIST := []string{"all", "validate_cases", "test_cases"}
    if !utils.Contains(LIST, mtype) {
        return fmt.Errorf("unsupport module type")
    }
    err = runTestlibValidators(&session.JudgeConfig, mtype, mCaseIndex)
    if err != nil {
        return err
    }
    return session.SaveConfiguration(true)
}