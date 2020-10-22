package test

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/client"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/LanceLRQ/deer-executor/provider"
	uuid "github.com/satori/go.uuid"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func initWorkRoot() error {
	_, filename, _ , _ := runtime.Caller(1)
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
	return nil
}

func runAPlusB(codeFile string, codeLang string) (*executor.JudgeResult, error) {
	session, err := executor.NewSession("./data/problems/APlusB/problem.json")
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

func analysisResult (caseName string, result *executor.JudgeResult, expect int) error {
	name, ok := executor.FlagMeansMap[result.JudgeResult]
	if !ok { name = "Unknown" }
	if result.JudgeResult != expect {
		ename, ok := executor.FlagMeansMap[expect]
		if !ok { ename = "Unknown" }
		return fmt.Errorf("[%s] expect %s, but got %s\n%s", caseName, ename, name, executor.ObjectToJSONStringFormatted(result))
	}
	fmt.Printf("[%s] finish with: %s\n", caseName, name)
	return nil
}

