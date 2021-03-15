package test

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/utils"
    "log"
    "math/rand"
    "runtime"
    "strconv"
    "testing"
    "time"
)

type CodeCases struct {
    Path string
    Expect int
}

func runJudgeBenchmark(codeFiles []CodeCases, b *testing.B) error {
    fileLen := len(codeFiles)
    rel := JudgementBenchmark{
        CurrectCounter: map[string]int{},
        IncurrectCounter: map[string]int{},
    }
    rand.Seed(time.Now().UnixNano())
    errFlag := false
    times := 100
    b.N = times

    fmt.Printf("%d tests.\n", times)

    startTime := time.Now().UnixNano()
    cExitCounter := map[int]int{}
    icExitCounter := map[int]int{}
    b.ResetTimer()
    for i := 0; i < times; i++ {

        j := rand.Int31n(int32(fileLen))
        fileCase := codeFiles[j]
        expect := fileCase.Expect

        if i % 10 == 0 {
            log.Printf("[%d / %d]\n", i, times)
        }
        judgeResult, err := runAPlusB(fileCase.Path, "")
        if err != nil {
            rel.Message = fmt.Sprintf("unexpected error! [%s]: %s\n", strconv.Itoa(i), err.Error())
            log.Printf(rel.Message)
            break
        }
        name, ok := constants.FlagMeansMap[judgeResult.JudgeResult]
        if !ok {
            name = "Unknown"
        }
        expName, ok := constants.FlagMeansMap[expect]
        if !ok {
            expName = "Unknown"
        }
        if judgeResult.JudgeResult != expect {
            log.Printf("unexpected! [%s] expect %s, got %s.\n", strconv.Itoa(i), expName, name)
            errFlag = true
            icExitCounter[judgeResult.JudgeResult]++
        } else {
            cExitCounter[judgeResult.JudgeResult]++
        }
    }
    b.StopTimer()
    endTime := time.Now().UnixNano()
    for key, value := range cExitCounter {
        name, ok := constants.FlagMeansMap[key]
        if !ok {
            name = "Unknown"
        }
        rel.CurrectCounter[name] = value
    }
    for key, value := range icExitCounter {
        name, ok := constants.FlagMeansMap[key]
        if !ok {
            name = "Unknown"
        }
        rel.IncurrectCounter[name] = value
    }
    rel.TimeUsed = float64(endTime - startTime) / float64(time.Second)

    fmt.Println(utils.ObjectToJSONStringFormatted(rel))

    if errFlag {
        b.Fatal("ERROR")
    }

    return nil
}

func BenchmarkAPlusBProblem(b *testing.B) {
    err := initWorkRoot()
    if err != nil {
        b.Fatal(err)
        return
    }
    mle := "./data/codes/APlusB/mle.c"
    if runtime.GOOS == "darwin" {
        mle = "./data/codes/APlusB/mle_darwin.c"
    }
    testCaseMap := []CodeCases {
        CodeCases{ Path: "./data/codes/APlusB/ac.c", Expect: constants.JudgeFlagAC },
        CodeCases{ Path: "./data/codes/APlusB/pe.c", Expect: constants.JudgeFlagPE },
        CodeCases{ Path: "./data/codes/APlusB/pe2.c", Expect: constants.JudgeFlagPE },
        CodeCases{ Path: "./data/codes/APlusB/pe3.c", Expect: constants.JudgeFlagPE },
        CodeCases{ Path: "./data/codes/APlusB/wa.c", Expect: constants.JudgeFlagWA },
        CodeCases{ Path: "./data/codes/APlusB/wa2.c", Expect: constants.JudgeFlagWA },
        CodeCases{ Path: mle, Expect: constants.JudgeFlagMLE },
        CodeCases{ Path: "./data/codes/APlusB/tle.c", Expect: constants.JudgeFlagTLE },
        CodeCases{ Path: "./data/codes/APlusB/ole.c", Expect: constants.JudgeFlagOLE },
        CodeCases{ Path: "./data/codes/APlusB/re.c", Expect: constants.JudgeFlagRE },
        CodeCases{ Path: "./data/codes/APlusB/re2.c", Expect: constants.JudgeFlagRE },
        CodeCases{ Path: "./data/codes/APlusB/ce.c", Expect: constants.JudgeFlagCE },
    }
    err = runJudgeBenchmark(testCaseMap, b)
    if err != nil {
        b.Fatal(err)
        return
    }
}