package packmgr

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "os"
)

func runValidator(validatorCase *structs.TestlibValidatorCase) error {
    return nil
}

func isValidatorExists (config *structs.JudgeConfiguration) error {
    validator, err := utils.GetCompiledBinaryFileAbsPath("validator", config.TestLib.ValidatorName, config.ConfigDir)
    if err != nil { return err }
    _, err = os.Stat(validator)
    if os.IsNotExist(err) {
        return fmt.Errorf("validator execuable binary not exists")
    }
    return nil
}

// 运行validator case的校验
// caseIndex < 0 表示校验全部
func RunTestlibValidatorCases(config *structs.JudgeConfiguration, caseIndex int) error {
    // 执行遍历
    if caseIndex < 0 {
        for _, vCase := range config.TestLib.ValidatorCases {
            err := runValidator(&vCase)
            if err != nil { return err }
        }
    } else {
        vCase := config.TestLib.ValidatorCases[caseIndex]
        err := runValidator(&vCase)
        if err != nil { return err }
    }
    return nil
}

// 运行Testlib的validator校验
func RunTestlibValidators(config *structs.JudgeConfiguration, moduleName string) error {
    if err := isValidatorExists(config); err != nil {
        return err
    }
    return nil
}
