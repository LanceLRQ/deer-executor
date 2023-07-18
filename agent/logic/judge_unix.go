//go:build linux || darwin
// +build linux darwin

package logic

import (
	"fmt"
	agentConfig "github.com/LanceLRQ/deer-executor/v3/agent/config"
	"github.com/LanceLRQ/deer-executor/v3/agent/rpc"
	"github.com/LanceLRQ/deer-executor/v3/executor"
	persistence "github.com/LanceLRQ/deer-executor/v3/executor/persistence"
	"github.com/LanceLRQ/deer-executor/v3/executor/persistence/result"
	commonStructs "github.com/LanceLRQ/deer-executor/v3/executor/structs"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// JudgementRunOption options for StartJudgement
type JudgementRunOption struct {
	Persistence *persistence.JudgeResultPersisOptions
	Clean       bool
	ConfigFile  string
	WorkDir     string
	ShowLog     bool
	LogLevel    int
	Language    string
	LibraryDir  string
	CodeStr     string
	SessionID   string
	SessionDir  string
	SessionRoot string
}

func checkJudgeRequsetArgs(request *rpc.JudgementRequest) error {
	if strings.TrimSpace(request.ProblemDir) == "" {
		return errors.Errorf("invalid problem path")
	}
	if strings.TrimSpace(request.Code) == "" {
		return errors.Errorf("invalid code file path")
	}
	return nil
}

func runOnceJudge(options *JudgementRunOption) (*commonStructs.JudgeResult, *executor.JudgeSession, error) {
	// create session
	session, err := executor.NewSessionWithLog(options.ConfigFile, options.ShowLog, options.LogLevel)
	if err != nil {
		return nil, nil, err
	}
	if options.Language != "" {
		session.CodeLangName = options.Language
	}
	if session.JudgeConfig.SpecialJudge.Mode > 0 {
		// check library when enable special judge
		libDir, err := filepath.Abs(options.LibraryDir)
		if err != nil {
			return nil, nil, errors.Errorf("get library root error: %s", err.Error())
		}
		s, err := os.Stat(libDir)
		if err != nil {
			return nil, nil, errors.Errorf("library root not exists")
		}
		if !s.IsDir() {
			return nil, nil, errors.Errorf("library root not a directory")
		}
		session.LibraryDir = libDir
	}
	// init files
	if options.WorkDir != "" {
		workDirAbsPath, err := filepath.Abs(options.WorkDir)
		if err != nil {
			return nil, nil, err
		}
		session.ConfigDir = workDirAbsPath
		session.JudgeConfig.ConfigDir = session.ConfigDir
	}
	session.CodeStr = options.CodeStr
	session.SessionID = options.SessionID
	session.SessionRoot = options.SessionRoot
	session.SessionDir = options.SessionDir
	// start judgement
	judgeResult := session.RunJudge()
	return &judgeResult, session, nil
}

func startRealJudgement(options *JudgementRunOption) (*executor.JudgeSession, *commonStructs.JudgeResult, error) {
	judgeResult, judgeSession, err := runOnceJudge(options)
	if err != nil {
		return nil, nil, err
	}
	// Do clean (or benchmark on)
	if options.Clean {
		defer judgeSession.Clean()
	}

	// persistence
	if options.Persistence != nil {
		options.Persistence.SessionDir = judgeSession.SessionDir
		err = result.PersistentJudgeResult(judgeResult, options.Persistence)
		if err != nil {
			return nil, nil, err
		}
	}
	return judgeSession, judgeResult, nil
}

// intro
func runRpcJudge(request *rpc.JudgementRequest) (*commonStructs.JudgeResult, string, error) {
	// build and check workdir
	workDir := path.Join(agentConfig.JudgementConfig.ProblemRoot, request.ProblemDir)
	if wd, err := os.Stat(workDir); os.IsNotExist(err) || (wd != nil && !wd.IsDir()) {
		return nil, "", fmt.Errorf("problem_dir is not exists or not a directory")
	}
	// build and check config file
	configFile := path.Join(workDir, "problem.json")
	if cf, err := os.Stat(configFile); os.IsNotExist(err) || (cf != nil && cf.IsDir()) {
		return nil, "", fmt.Errorf("config_file is not exists or not a file")
	}
	// build persistence options
	compressorType := uint8(request.CompressType)
	jOption := persistence.JudgeResultPersisOptions{
		CompressorType:   compressorType,
		SaveAcceptedData: request.PersistWithAcData,
	}
	// Is enable persistence with sign
	if request.PersistResult && request.SignResult {
		passphrase := []byte(request.GpgPassphrase)
		pem, err := persistence.GetArmorPublicKey(request.GpgKey, passphrase)
		if err != nil {
			return nil, "", err
		}
		jOption.DigitalSign = true
		jOption.DigitalPEM = pem
	}

	// create session info
	sessionID := uuid.NewV1().String()
	// init session dir
	sessionDir, err := utils.GetSessionDir(agentConfig.JudgementConfig.SessionRoot, sessionID)
	if err != nil {
		return nil, "", err
	}

	// set log and log level
	showLog := request.EnableLog
	logLevel := int(request.LogLevel)

	// build options
	rOptions := &JudgementRunOption{
		Clean:       request.CleanSession,
		ShowLog:     showLog,
		LogLevel:    logLevel,
		WorkDir:     workDir,
		ConfigFile:  configFile,
		Language:    request.Language,
		LibraryDir:  agentConfig.JudgementConfig.SystemLibraryRoot,
		CodeStr:     request.Code,
		SessionID:   sessionID,
		SessionDir:  sessionDir,
		SessionRoot: agentConfig.JudgementConfig.SessionRoot,
	}

	persistFile := ""          // set empty
	if request.PersistResult { // if enable persistence
		persistFile = fmt.Sprintf("%s.result", sessionID)
		// outfile
		jOption.OutFile = path.Join(agentConfig.JudgementConfig.SessionRoot, persistFile)
		rOptions.Persistence = &jOption
	}

	// start judge
	_, judgeResult, err := startRealJudgement(rOptions)
	if err != nil {
		return nil, "", err
	}

	return judgeResult, persistFile, nil
}
