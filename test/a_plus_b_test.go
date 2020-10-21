package test

import (
	"github.com/LanceLRQ/deer-executor/executor"
	"testing"
)


// Test 1: AC
func TestAPlusBProblemAc(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ac.c")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	err = analysisResult(result, executor.JudgeFlagAC)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Log("OK")
}

// Test 2: PE which no space char
func TestAPlusBProblemPE1(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/pe.c")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	err = analysisResult(result, executor.JudgeFlagPE)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Log("OK")
}

// Test 3: PE which out-of-order space char
func TestAPlusBProblemPE2(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/pe2.c")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	err = analysisResult(result, executor.JudgeFlagPE)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Log("OK")
}

// Test 4: PE which wrong space char
func TestAPlusBProblemPE3(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/pe3.c")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	err = analysisResult(result, executor.JudgeFlagPE)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Log("OK")
}

// Test 5: WA because of using multiple
func TestAPlusBProblemWA(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/wa.c")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult(result, executor.JudgeFlagWA)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

// Test 6: WA because of reading only one line
func TestAPlusBProblemWA2(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/wa2.c")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult(result, executor.JudgeFlagWA)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

// Test 7: use Java language
func TestAPlusBProblemJava(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ac.java")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult(result, executor.JudgeFlagAC)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

// Test 8: use Go language
func TestAPlusBProblemGo(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ac.go")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult(result, executor.JudgeFlagAC)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}