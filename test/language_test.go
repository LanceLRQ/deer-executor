package test

import (
	"flag"
	"github.com/LanceLRQ/deer-executor/v3/executor/constants"
	"testing"
)

// Test: CPP language
func TestAPlusBProblemCpp(t *testing.T) {
	if flag.Arg(0) != "all-language" {
		t.Log("CPP: Skip")
		return
	}
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ac.cpp", "")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("test cpp", result, constants.JudgeFlagAC)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

// Test: Java language
func TestAPlusBProblemJava(t *testing.T) {
	if flag.Arg(0) != "all-language" {
		t.Log("Java: Skip")
		return
	}
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ac.java", "")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("test java", result, constants.JudgeFlagAC)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

//// Test: Python 2 language
//func TestAPlusBProblemPython2(t *testing.T) {
//	if flag.Arg(0) != "all-language" {
//		t.Log("Python 2: Skip")
//		return
//	}
//	err := initWorkRoot()
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	result, err := runAPlusB("./data/codes/APlusB/ac_py2.py", "python2")
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	err = analysisResult("test python 2", result, executor.JudgeFlagAC)
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	t.Log("OK")
//}

// Test: Python 3 language
func TestAPlusBProblemPython3(t *testing.T) {
	if flag.Arg(0) != "all-language" {
		t.Log("Python 3: Skip")
		return
	}
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ac_py3.py", "python3")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("test python 3", result, constants.JudgeFlagAC)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

// Test: Python 3 language ce - indent error
func TestAPlusBProblemPython3CE(t *testing.T) {
	if flag.Arg(0) != "all-language" {
		t.Log("Python 3 CE 1: Skip")
		return
	}
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ce_py3.py", "python3")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("test python 3 ce 1", result, constants.JudgeFlagCE)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

// Test: Python 3 language ce - syntax error
func TestAPlusBProblemPython3CE2(t *testing.T) {
	if flag.Arg(0) != "all-language" {
		t.Log("Python 3 CE 2: Skip")
		return
	}
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ce2_py3.py", "python3")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("test python 3 ce 2", result, constants.JudgeFlagCE)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

// Test: Go language
func TestAPlusBProblemGo(t *testing.T) {
	if flag.Arg(0) != "all-language" {
		t.Log("Go: Skip")
		return
	}
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ac.go", "")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("test golang", result, constants.JudgeFlagAC)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

//// Test: Nodejs language
//func TestAPlusBProblemNode(t *testing.T) {
//	if flag.Arg(0) != "all-language" {
//		t.Log("NodeJS: Skip")
//		return
//	}
//	err := initWorkRoot()
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	result, err := runAPlusB("./data/codes/APlusB/ac.js")
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	err = analysisResult("test nodejs", result, executor.JudgeFlagAC)
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	t.Log("OK")
//}

// Test: PHP language
func TestAPlusBProblemPHP(t *testing.T) {
	if flag.Arg(0) != "all-language" {
		t.Log("PHP: Skip")
		return
	}
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runAPlusB("./data/codes/APlusB/ac.php", "")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("test php", result, constants.JudgeFlagAC)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

//
//// Test: Ruby language
//func TestAPlusBProblemRuby(t *testing.T) {
//	if flag.Arg(0) != "all-language" {
//		t.Log("Ruby: Skip")
//		return
//	}
//	err := initWorkRoot()
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	result, err := runAPlusB("./data/codes/APlusB/ac.rb")
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	err = analysisResult("test ruby", result, executor.JudgeFlagAC)
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	t.Log("OK")
//}
