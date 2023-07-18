package test

import (
	"github.com/LanceLRQ/deer-executor/v3/executor/constants"
	"testing"
)

// Test: AC
func TestWJ2018AC(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runWJ2018("./data/codes/WJ2018/answer_ac.c", "")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("wj2018-1", result, constants.JudgeFlagAC)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

// Test: WA 1
func TestWJ2018WA1(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runWJ2018("./data/codes/WJ2018/answer_wa.c", "")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("wj2018-2", result, constants.JudgeFlagWA)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

// Test: WA 2
func TestWJ2018WA2(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runWJ2018("./data/codes/WJ2018/answer_wa2.c", "")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("wj2018-3", result, constants.JudgeFlagWA)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}
