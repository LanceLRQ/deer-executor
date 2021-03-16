// +build linux darwin

package run

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/v2/client"
    "github.com/urfave/cli/v2"
    "log"
    "os"
    "strconv"
    "time"
)

type JudgementBenchmark struct {
    TimeUsed float64        `json:"time_used"`
    Counter  map[string]int `json:"counter"`
    Message  string         `json:"message"`
}

func runJudgeBenchmark(c *cli.Context, configFile string) error {
    rel := JudgementBenchmark{
        Counter: map[string]int{},
    }

    times := c.Int("benchmark")
    rfp, err := os.OpenFile("./report.log", os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer rfp.Close()

    workDir := c.String("work-dir")
    // 构建运行选项
    rOptions := &RunOption{
        Clean: true,
        ShowLog: false,
        LogLevel: 0,
        WorkDir: workDir,
        ConfigFile: configFile,
        Language: c.String("language"),
        LibraryDir: c.String("library"),
        CodePath: c.Args().Get(1),
        SessionId: "",
        SessionRoot: "",
    }

    startTime := time.Now().UnixNano()
    exitCounter := map[int]int{}
    for i := 0; i < times; i++ {
        if i % 10 == 0 {
            log.Printf("[%d / %d]\n", i, times)
        }
        judgeResult, _, err := runOnceJudge(rOptions)
        if err != nil {
            rel.Message = fmt.Sprintf("break! %s\n", err.Error())
            log.Print(rel.Message)
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
        rel.Counter[name] = value
        log.Printf("%s: %d\n", name, value)
    }
    rel.TimeUsed = float64(endTime - startTime)
    duration := rel.TimeUsed / float64(time.Second)
    log.Printf("total time used: %.2fs\n", duration)

    client.NewClientSuccessMessage(rel).Print(true)
    return nil
}
