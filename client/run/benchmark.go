package run

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/urfave/cli/v2"
    "os"
    "strconv"
    "time"
)

func runJudgeBenchmark (c *cli.Context, configFile string) error {
    times := c.Int("benchmark")
    rfp, err := os.OpenFile("./report.log", os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer rfp.Close()

    startTime := time.Now().UnixNano()
    exitCounter := map[int]int{}
    for i := 0; i < times; i++ {
        if i%10 == 0 {
            fmt.Printf("%d / %d\n", i, times)
        }
        judgeResult, _, err := runOnceJudge(c, configFile, i)
        if err != nil {
            fmt.Printf("break! %s\n", err.Error())
            _, _ = rfp.WriteString(fmt.Sprintf("[%s]: %s\n", strconv.Itoa(i), err.Error()))
            break
        }
        name, ok := constants.FlagMeansMap[judgeResult.JudgeResult]
        if !ok {
            name = "Unknown"
        }
        _, _ = rfp.WriteString(fmt.Sprintf("[%s]: %s\n", judgeResult.SessionId, name))
        if judgeResult.JudgeResult != constants.JudgeFlagAC {
            _, _ = rfp.WriteString(utils.ObjectToJSONStringFormatted(judgeResult) + "\n")
        }
        exitCounter[judgeResult.JudgeResult]++
    }
    endTime := time.Now().UnixNano()
    for key, value := range exitCounter {
        name, ok := constants.FlagMeansMap[key]
        if !ok {
            name = "Unknown"
        }
        fmt.Printf("%s: %d\n", name, value)
    }
    duration := float64(endTime-startTime) / float64(time.Second)
    fmt.Printf("total time used: %.2fs\n", duration)
    return nil
}
