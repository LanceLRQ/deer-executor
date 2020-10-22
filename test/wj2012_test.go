package test

import (
	"github.com/LanceLRQ/deer-executor/executor"
	"testing"
)

// Test: AC
func TestWJ2012AC(t *testing.T) {
	err := initWorkRoot()
	if err != nil {
		t.Fatal(err)
		return
	}
	result, err := runWJ2012("./data/codes/WJ2012/answer_ac.c", "")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = analysisResult("wj2012", result, executor.JudgeFlagAC)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("OK")
}

