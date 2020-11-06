package run

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/persistence/problems"
    "github.com/LanceLRQ/deer-common/provider"
    uuid "github.com/satori/go.uuid"
    "os"
)

func loadSystemConfiguration () error {
    // 载入默认配置
    err := provider.PlaceCompilerCommands("./compilers.json")
    if err != nil {
        return err
    }
    err = constants.PlaceMemorySizeForJIT("./jit_memory.json")
    if err != nil {
        return err
    }
    return nil
}

func loadProblemConfiguration(configFile string, workDir string) (string, bool, string, error) {
    _, err := os.Stat(configFile)
    if err != nil && os.IsNotExist(err) {
        return "", false, "", fmt.Errorf("problem config file (%s) not found", configFile)
    }
    yes, err := problems.IsProblemPackage(configFile)
    if err != nil {
        return "", false, "", err
    }
    autoRemoveWorkDir := false
    // 如果是题目包文件，进行解包
    if yes {
        if workDir == "" {
            workDir = "/tmp/" + uuid.NewV4().String()
            autoRemoveWorkDir = true
        }
        if info, err := os.Stat(workDir); os.IsNotExist(err) {
            err = os.MkdirAll(workDir, 0755)
            if err != nil {
                return "", false, "", err
            }
        } else if !info.IsDir() {
            return "", false, "", fmt.Errorf("work dir path cannot be a file path")
        }
        _, newConfigFile, err := problems.ReadProblemInfo(configFile, true, workDir)
        if err != nil {
            return "", false, "", err
        }
        configFile = newConfigFile
    }
    return configFile, autoRemoveWorkDir, workDir, nil
}