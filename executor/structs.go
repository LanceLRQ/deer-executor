package executor

import (
	"encoding/json"
	"fmt"
	"github.com/LanceLRQ/deer-executor/v2/common/logger"
	"github.com/LanceLRQ/deer-executor/v2/common/provider"
	commonStructs "github.com/LanceLRQ/deer-executor/v2/common/structs"
	"github.com/LanceLRQ/deer-executor/v2/common/utils"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// JudgeSession 评测会话类
type JudgeSession struct {
	SessionID    string   // Judge Session Id
	SessionRoot  string   // Session Root Directory
	SessionDir   string   // Session Directory
	ConfigFile   string   // Config file
	ConfigDir    string   // Config file dir
	CodeLangName string   // Code file language name
	CodeFile     string   // Code File Path
	CodeStr      string   // Code Str (if set, use it first)
	LibraryDir   string   // Compile Library Path for Working Program
	Commands     []string // Executable program commands

	JudgeConfig commonStructs.JudgeConfiguration      // Judge Configurations
	Compiler    provider.CodeCompileProviderInterface // Compiler entity

	Logger  *logger.JudgeLogger // Judge Logger
	Timeout int                 // Process timeout (s)
}

// SaveConfiguration 保存评测会话
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

// NewSession 创建会话对象
func NewSession(configFile string) (*JudgeSession, error) {
	session := JudgeSession{}
	session.Logger = logger.NewJudgeLogger()
	session.Timeout = 30 // 默认30秒超时
	session.SessionRoot = "/tmp"
	session.CodeLangName = "auto"
	session.JudgeConfig.UID = -1
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

// NewSessionWithLog 创建一个包含日志的对象
func NewSessionWithLog(configFile string, print bool, level int) (*JudgeSession, error) {
	session, err := NewSession(configFile)
	if err != nil {
		return nil, err
	}
	session.Logger.SetStdoutPrint(print)
	session.Logger.SetLogLevel(level)
	return session, nil
}

// Clean 清理会话工作目录
func (session *JudgeSession) Clean() {
	_ = os.RemoveAll(session.SessionDir)
}
