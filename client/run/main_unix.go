// +build linux darwin

package run

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/urfave/cli/v2"
    "os"
)

// 执行评测
func UserRunJudge(c *cli.Context) error {
    err := loadSystemConfiguration()
    if err != nil {
        return err
    }

    configFile, autoRemoveWorkDir, workDir, err := loadProblemConfiguration(c.Args().Get(0), c.String("work-dir"))
    if err != nil {
        return err
    }
    if autoRemoveWorkDir {
        defer (func() {
            _ = os.RemoveAll(workDir)
        })()
    }

    isBenchmarkMode := c.Int("benchmark") > 1
    if !isBenchmarkMode {
        // 普通的运行
        judgeResult, err := runUserJudge(c, configFile, workDir)
        if err != nil {
            return err
        }
        fmt.Println(utils.ObjectToJSONStringFormatted(judgeResult))
        os.Exit(judgeResult.JudgeResult)
    } else {
        // 基准测试
        err = runJudgeBenchmark(c, configFile)
        if err != nil {
            return err
        }
    }
    return nil
}
