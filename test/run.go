package main

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor"
	"github.com/LanceLRQ/deer-executor/compile"
	"io/ioutil"
)

func main() {

	workDir := "/Users/lancelrq/wejudge/deer-executor"

	//code, err:= ioutil.ReadFile(workDir + "/test/program.c")
	code, err:= ioutil.ReadFile(workDir + "/test/Main.java")
	if err != nil {
		return
	}
	codeContent := string(code)

	//compiler := deer_compiler.GnucCompileProvider{}
	compiler := deer_compiler.JavaCompileProvider{}
	compiler.Init(codeContent,  "/tmp")
	success, ceinfo := compiler.Compile()
	if !success {
		fmt.Println(ceinfo)
	}

	judgeOption := deer_executor.JudgeOption{
		TimeLimit:     1000,
		MemoryLimit:   256000,
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
	judgeResult, err := deer_executor.Judge(judgeOption)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Println(judgeResult.String())
	}
}