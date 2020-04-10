package test

import (
	deer_executor "github.com/LanceLRQ/deer-executor"
	deer_compiler "github.com/LanceLRQ/deer-executor/compile"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func SpecialJudge (workDir, codeContent, handle string, t *testing.T) *deer_executor.JudgeResult {

	compiler := deer_compiler.GnucCompileProvider{}
	//compiler := deer_compiler.JavaCompileProvider{}
	err := compiler.Init(codeContent,  "/tmp")
	if err != nil {
		log.Fatal(err)
	}
	success, ceinfo := compiler.Compile()
	if !success {
		t.Fatal(ceinfo)
	}

	spjcode, err := ioutil.ReadFile(workDir + "/scripts/special_judger.cpp")
	if err != nil {
		t.Fatal(err.Error())
	}
	spjCodeContent := string(spjcode)
	spjCompiler := deer_compiler.GnucppCompileProvider{}
	err = spjCompiler.Init(spjCodeContent,  "/tmp")
	if err != nil {
		log.Fatal(err)
	}
	success, ceinfo = spjCompiler.Compile()
	if !success {
		t.Fatal(ceinfo)
	}

	_ = os.Remove("/tmp/user.out")
	_ = os.Remove("/tmp/user.err")
	_ = os.Remove("/tmp/spj.out")
	_ = os.Remove("/tmp/spj.err")

	judgeOption := deer_executor.JudgeOption{
		TimeLimit:     10000,
		MemoryLimit:   65355,
		FileSizeLimit: 100 * 1024 * 1024,
		Commands:      compiler.GetRunArgs(),
		TestCaseIn:    workDir + "/cases/" + handle + ".in",
		TestCaseOut:   workDir + "/cases/" + handle + ".out",
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
			Mode:        deer_executor.SPECIAL_JUDGE_MODE_CHECKER,
			Checker:     spjCompiler.GetRunArgs()[0],
			RedirectStd: true,
			TimeLimit:   10000,
			MemoryLimit: 65535,
			Stdout:      "/tmp/spj.out",
			Stderr:      "/tmp/spj.err",
		},
		Uid: -1,
	}
	judgeResult, err := deer_executor.Judge(judgeOption)
	if err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(judgeResult.String())
	}
	return judgeResult
}

func TestSpecialJudgeAC(t *testing.T) {
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	code, err:= ioutil.ReadFile(workDir + "/scripts/special_answer.c")
	if err != nil {
		t.Fatal(err.Error())
	}
	codeContent := string(code)
	rel := SpecialJudge(workDir, codeContent, "0", t)
	if rel.JudgeResult != deer_executor.JUDGE_FLAG_AC {
		t.Fatal("Program not AC")
	} else {
		t.Log("OK")
	}
}
