package main

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/src"
	"github.com/LanceLRQ/deer-executor/src/compile"
	"io/ioutil"
)

func main() {

	workDir := "/Users/lancelrq/wejudge/deer-executor"

	code, err:= ioutil.ReadFile(workDir + "/test/program.c")
	if err != nil {
		return
	}
	codeContent := string(code)

	compiler := compile.GnucCompileProvider{}
	compiler.Init(codeContent,  "/tmp")
	success, ceinfo := compiler.Compile()
	if !success {
		fmt.Println(ceinfo)
	}

	judgeOption := deer.JudgeOption{
		TimeLimit:     1000,
		MemoryLimit:   32768,
		FileSizeLimit: 100 * 1024 * 1024,
		Commands:      compiler.GetRunArgs(),
		TestCaseIn:    workDir + "/test/0.in",
		TestCaseOut:   workDir + "/test/0.out",
		ProgramOut:    "/tmp/user.out",
		ProgramError:  "/tmp/user.err",
		// Special Judge
		SpecialJudge: struct {
			Mode        int
			Checker     string
			RedirectStd bool
			TimeLimit   int
			MemoryLimit int
			Stdout      string
			Stderr      string
		}{
			Mode:        0,
			Checker:     "",
			RedirectStd: false,
			TimeLimit:   0,
			MemoryLimit: 0,
			Stdout:      "",
			Stderr:      "",
		},
		Uid: -1,
	}
	judgeResult, err := deer.Judge(judgeOption)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Println(judgeResult.String())
	}
}