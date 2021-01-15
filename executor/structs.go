package executor

import (
    "encoding/json"
    "fmt"
    "github.com/LanceLRQ/deer-common/logger"
    "github.com/LanceLRQ/deer-common/provider"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "strings"
)

// 评测会话类
type JudgeSession struct {
    SessionId    string   // Judge Session Id
    SessionRoot  string   // Session Root Directory
    SessionDir   string   // Session Directory
    ConfigFile   string   // Config file
    ConfigDir    string   // Config file dir
    CodeLangName string   // Code file language name
    CodeFile     string   // Code File Path
    LibraryDir   string   // Compile Library Path for Working Program
    Commands     []string // Executable program commands

    JudgeConfig commonStructs.JudgeConfiguration      // Judge Configurations
    Compiler    provider.CodeCompileProviderInterface // Compiler entity

    Logger      *logger.JudgeLogger // Judge Logger
}

// 保存评测会话
func (session *JudgeSession) SaveConfiguration(userConfirm bool) error {
    if userConfirm {
        fmt.Print("Save all the changed to config file? [y/N] ")
        ans := ""
        _, err := fmt.Scanf("%s", &ans)
        if err != nil {
            return nil // don't crash at EOF
        }
        if len(ans) > 0 && strings.ToLower(ans[:1]) != "y" {
            return nil
        }
    }
    err := ioutil.WriteFile(session.ConfigFile, []byte(utils.ObjectToJSONStringFormatted(session.JudgeConfig)), 0644)
    if err != nil {
        return err
    }
    if userConfirm {
        fmt.Println("Saved!")
    }
    return nil
}

// 创建会话对象
func NewSession(configFile string) (*JudgeSession, error) {
    session := JudgeSession{}
    session.Logger = logger.NewJudgeLogger()
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
        session.JudgeConfig.ConfigDir = session.ConfigDir
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

// 创建一个包含日志的对象
func NewSessionWithLog(configFile string, print bool, level int) (*JudgeSession, error) {
    session, err := NewSession(configFile)
    if err != nil {
        return nil, err
    }
    session.Logger.SetStdoutPrint(print)
    session.Logger.SetLogLevel(level)
    return session, nil
}

// 清理案发现场
func (session *JudgeSession) Clean() {
    _ = os.RemoveAll(session.SessionDir)
}
