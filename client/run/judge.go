// +build linux darwin

package run

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/logger"
    "github.com/LanceLRQ/deer-common/persistence"
    "github.com/LanceLRQ/deer-common/persistence/judge_result"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/executor"
    uuid "github.com/satori/go.uuid"
    "github.com/urfave/cli/v2"
    "log"
    "os"
    "path/filepath"
    "strconv"
)

// 执行一次完整的评测
func runOnceJudge(c *cli.Context, configFile, workDir string, counter int) (*commonStructs.JudgeResult, *executor.JudgeSession, error) {
    isBenchmarkMode := c.Int("benchmark") > 1
    var logLevel int
    showLog := false
    if c.Bool("log") {
        showLog = true
        var ok bool
        logLevelStr := c.String("log-level")
        logLevel, ok = logger.LogLevelStrMapping[logLevelStr]
        if !ok {
            logLevel = 0
        }
    }
    // create session
    session, err := executor.NewSessionWithLog(configFile, !isBenchmarkMode && showLog, logLevel)
    if err != nil {
        return nil, nil, err
    }
    if c.String("language") != "" {
        session.CodeLangName = c.String("language")
    }
    if session.JudgeConfig.SpecialJudge.Mode > 0 {
        // 特判时需要检查library目录
        libDir, err := filepath.Abs(c.String("library"))
        if err != nil {
            return nil, nil, fmt.Errorf("get library root error: %s", err.Error())
        }
        if s, err := os.Stat(libDir); err != nil {
            return nil, nil, fmt.Errorf("library root not exists")
        } else {
            if !s.IsDir() {
                return nil, nil, fmt.Errorf("library root not a directory")
            }
        }
        session.LibraryDir = libDir
    }
    // init files
    if workDir != "" {
        workDirAbsPath, err := filepath.Abs(workDir)
        if err != nil {
            return nil, nil, err
        }
        session.ConfigDir = workDirAbsPath
        session.JudgeConfig.ConfigDir = session.ConfigDir
    }
    session.CodeFile = c.Args().Get(1)
    session.SessionId = c.String("session-id")
    session.SessionRoot = c.String("session-root")
    // create session info
    if isBenchmarkMode {
        session.SessionId = uuid.NewV1().String() + strconv.Itoa(counter)
    } else {
        if session.SessionId == "" {
            session.SessionId = uuid.NewV1().String()
        }
    }
    if session.SessionRoot == "" {
        session.SessionRoot = "/tmp"
    }
    // 初始化session dir
    sessionDir, err := utils.GetSessionDir(session.SessionRoot, session.SessionId)
    if err != nil {
        log.Fatal(err)
        return nil, nil, err
    }
    session.SessionDir = sessionDir
    // start judgement
    judgeResult := session.RunJudge()
    return &judgeResult, session, nil
}

func runUserJudge (c *cli.Context, configFile, workDir string) (*commonStructs.JudgeResult, error) {
    // parse params
    persistenceOn := c.String("persistence") != ""
    digitalSign := c.Bool("sign")
    compressorType := uint8(1)
    if c.String("compress") == "none" {
        compressorType = uint8(0)
    }
    jOption := persistence.JudgeResultPersisOptions{
        CompressorType: compressorType,
        SaveAcceptedData: c.Bool("save-ac-data"),
    }
    jOption.OutFile = c.String("persistence")
    // 是否要持久化结果
    if persistenceOn {
        if digitalSign {
            if c.String("passphrase") != "" {
                log.Println("[warn] Using a password on the command line interface can be insecure.")
            }
            passphrase := []byte(c.String("passphrase"))
            pem, err := persistence.GetArmorPublicKey(c.String("gpg-key"), passphrase)
            if err != nil {
                return nil, err
            }
            jOption.DigitalSign = true
            jOption.DigitalPEM = pem
        }
    }
    // Start Judgement
    judgeResult, judgeSession, err := runOnceJudge(c, configFile, workDir, 0)
    if err != nil {
        return nil, err
    }
    // Do clean (or benchmark on)
    if !c.Bool("no-clean") {
        defer judgeSession.Clean()
    }

    // persistence
    jOption.SessionDir = judgeSession.SessionDir
    if persistenceOn {
        err = judge_result.PersistentJudgeResult(judgeResult, &jOption)
        if err != nil {
            return nil, err
        }
    }
    if !c.Bool("detail") {
        judgeResult.TestCases = nil
        judgeResult.JudgeLogs = nil
    }
    return judgeResult, nil
}