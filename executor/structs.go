package executor

import (
    "encoding/json"
    "github.com/LanceLRQ/deer-common/provider"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "io/ioutil"
    "path"
    "path/filepath"
    "syscall"
)

type ProcessInfo struct {
    Pid    uintptr            `json:"pid"`
    Status syscall.WaitStatus `json:"status"`
    Rusage syscall.Rusage     `json:"rusage"`
}

// Judge session
type JudgeSession struct {
    SessionId    string   // Judge Session Id
    SessionRoot  string   // Session Root Directory
    SessionDir   string   // Session Directory
    ConfigFile   string   // Config file
    ConfigDir    string   // Config file dir
    CodeLangName string   // Code file language name
    CodeFile     string   // Code File Path
    Commands     []string // Executable program commands

    JudgeConfig commonStructs.JudgeConfiguration // Judge Configurations

    Compiler provider.CodeCompileProviderInterface // Compiler entity
}

func NewSession(configFile string) (*JudgeSession, error) {
    session := JudgeSession{}
    session.SessionRoot = "/tmp"
    session.CodeLangName = "auto"
    session.JudgeConfig.Uid = -1
    session.JudgeConfig.TimeLimit = 1000
    session.JudgeConfig.MemoryLimit = 65535
    session.JudgeConfig.StrictMode = true
    session.JudgeConfig.FileSizeLimit = 50 * 1024 * 1024
    session.JudgeConfig.SpecialJudge.Mode = 0
    session.JudgeConfig.SpecialJudge.RedirectProgramOut = true
    session.JudgeConfig.SpecialJudge.TimeLimit = 1000
    session.JudgeConfig.SpecialJudge.MemoryLimit = 65535
    if configFile != "" {
        configFileAbsPath, err := filepath.Abs(configFile)
        if err != nil {
            return nil, err
        }
        session.ConfigFile = configFileAbsPath
        session.ConfigDir = path.Dir(configFileAbsPath)
        cbody, err := ioutil.ReadFile(configFileAbsPath)
        if err != nil {
            return nil, err
        }
        err = json.Unmarshal(cbody, &session.JudgeConfig)
        if err != nil {
            return nil, err
        }
    }
    return &session, nil
}

func (*JudgeSession) GetCompiledBinary() {}