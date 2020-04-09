package test

import (
	deer_executor "github.com/LanceLRQ/deer-executor"
	deer_compiler "github.com/LanceLRQ/deer-executor/compile"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Runner (workDir, codeContent, handle string, t *testing.T) *deer_executor.JudgeResult {

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

	_ = os.Remove("/tmp/user.out")
	_ = os.Remove("/tmp/user.err")

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
		t.Fatal(err.Error())
	} else {
		t.Log(judgeResult.String())
	}
	return judgeResult
}

func TestNormalRunnerAC(t *testing.T) {
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	code, err:= ioutil.ReadFile(workDir + "/scripts/ac.c")
	if err != nil {
		t.Fatal(err.Error())
	}
	codeContent := string(code)
	rel := Runner(workDir, codeContent, "0", t)
	if rel.JudgeResult != deer_executor.JUDGE_FLAG_AC {
		t.Fatal("Program not AC")
	} else {
		t.Log("OK")
	}
}
func TestNormalRunnerPE(t *testing.T) {
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	code, err:= ioutil.ReadFile(workDir + "/scripts/pe.c")
	if err != nil {
		t.Fatal(err.Error())
	}
	codeContent := string(code)
	rel := Runner(workDir, codeContent, "0", t)
	if rel.JudgeResult != deer_executor.JUDGE_FLAG_PE {
		t.Fatal("Program not PE")
	} else {
		t.Log("OK")
	}
}
func TestNormalRunnerPETab(t *testing.T) {
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	code, err:= ioutil.ReadFile(workDir + "/scripts/pe2.c")
	if err != nil {
		t.Fatal(err.Error())
	}
	codeContent := string(code)
	rel := Runner(workDir, codeContent, "0", t)
	if rel.JudgeResult != deer_executor.JUDGE_FLAG_PE {
		t.Fatal("Program not PE")
	} else {
		t.Log("OK")
	}
}
func TestNormalRunnerWA(t *testing.T) {
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	code, err:= ioutil.ReadFile(workDir + "/scripts/wa.c")
	if err != nil {
		t.Fatal(err.Error())
	}
	codeContent := string(code)
	rel := Runner(workDir, codeContent,"0", t)
	if rel.JudgeResult != deer_executor.JUDGE_FLAG_WA {
		t.Fatal("Program not WA")
	} else {
		t.Log("OK")
	}
}