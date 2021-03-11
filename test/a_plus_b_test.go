package test

import (
    "github.com/LanceLRQ/deer-common/constants"
    "runtime"
    "testing"
)

// Test: AC
func TestAPlusBProblemAc(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/ac.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 1", result, constants.JudgeFlagAC)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: PE which no space char
func TestAPlusBProblemPE1(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/pe.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 2", result, constants.JudgeFlagPE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: PE which out-of-order space char
func TestAPlusBProblemPE2(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/pe2.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 3", result, constants.JudgeFlagPE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: PE which wrong space char
func TestAPlusBProblemPE3(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/pe3.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 4", result, constants.JudgeFlagPE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: Compile error
func TestAPlusBProblemCE(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/ce.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 5", result, constants.JudgeFlagCE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: memory out of limit
func TestAPlusBProblemMLE(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    p := "./data/codes/APlusB/mle.c"
    if runtime.GOOS == "darwin" {
        p = "./data/codes/APlusB/mle_darwin.c"
    }
    result, err := runAPlusB(p, "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 6", result, constants.JudgeFlagMLE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: output content out of limit (> 50m)
func TestAPlusBProblemOLE(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/ole.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 7", result, constants.JudgeFlagOLE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: output content out of limit (double answer)
func TestAPlusBProblemOLE2(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/ole2.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 8", result, constants.JudgeFlagOLE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: Runtime error (zero divide)
func TestAPlusBProblemRE(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/re.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 9", result, constants.JudgeFlagRE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: Runtime error (array out of bounds)
func TestAPlusBProblemRE2(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/re2.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 10", result, constants.JudgeFlagRE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: time out of limit (Endless loop)
func TestAPlusBProblemTLE(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/tle.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 11", result, constants.JudgeFlagTLE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: time out of limit (sleep)
func TestAPlusBProblemTLE2(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/tle2.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 12", result, constants.JudgeFlagTLE)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: WA because of using multiple
func TestAPlusBProblemWA(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/wa.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 13", result, constants.JudgeFlagWA)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}

// Test: WA because of reading only one line
func TestAPlusBProblemWA2(t *testing.T) {
    err := initWorkRoot()
    if err != nil {
        t.Fatal(err)
        return
    }
    result, err := runAPlusB("./data/codes/APlusB/wa2.c", "")
    if err != nil {
        t.Fatal(err)
        return
    }
    err = analysisResult("case 14", result, constants.JudgeFlagWA)
    if err != nil {
        t.Fatal(err)
        return
    }
    t.Log("OK")
}
