package test

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/provider"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/client"
    "github.com/LanceLRQ/deer-executor/executor"
    uuid "github.com/satori/go.uuid"
    "os"
    "path"
    "path/filepath"
    "runtime"
)

func initWorkRoot() error {
    _, filename, _, _ := runtime.Caller(1)
    workPath, err := filepath.Abs(path.Dir(path.Dir(filename)))
    if err != nil {
        return err
    }
    err = os.Chdir(workPath)
    if err != nil {
        return err
    }
    err = provider.PlaceCompilerCommands("./compilers.json")
    if err != nil {
        return err
    }
    err = constants.PlaceMemorySizeForJIT("./jit_memory.json")
    if err != nil {
        return err
    }
    return nil
}

func runJudge(conf, codeFile, codeLang string) (*commonStructs.JudgeResult, error) {
    session, err := executor.NewSession(conf)
    if err != nil {
        return nil, err
    }
    session.CodeFile = codeFile
    session.CodeLangName = codeLang
    session.SessionRoot = "/tmp"
    session.SessionId = uuid.NewV1().String()
    sessionDir, err := client.GetSessionDir(session.SessionRoot, session.SessionId)
    if err != nil {
        return nil, err
    }
    session.SessionDir = sessionDir
    defer session.Clean()
    // start judge
    judgeResult := session.RunJudge()
    return &judgeResult, err
}

func runAPlusB(codeFile, codeLang string) (*commonStructs.JudgeResult, error) {
    return runJudge("./data/problems/APlusB/problem.json", codeFile, codeLang)
}

func runWJ2018(codeFile, codeLang string) (*commonStructs.JudgeResult, error) {
    return runJudge("./data/problems/WJ2018/problem.json", codeFile, codeLang)
}

func runWJ2012(codeFile, codeLang string) (*commonStructs.JudgeResult, error) {
    return runJudge("./data/problems/WJ2012/problem.json", codeFile, codeLang)
}

func analysisResult(caseName string, result *commonStructs.JudgeResult, expect int) error {
    name, ok := constants.FlagMeansMap[result.JudgeResult]
    if !ok {
        name = "Unknown"
    }
    if result.JudgeResult != expect {
        ename, ok := constants.FlagMeansMap[expect]
        if !ok {
            ename = "Unknown"
        }
        return fmt.Errorf("[%s] expect %s, but got %s\n%s", caseName, ename, name, utils.ObjectToJSONStringFormatted(result))
    }
    fmt.Printf("[%s] finish with: %s\n", caseName, name)
    return nil
}
