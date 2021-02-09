// +build linux darwin

package run

import (
    "github.com/LanceLRQ/deer-executor/v2/client"
    "github.com/pkg/errors"
    "github.com/urfave/cli/v2"
    "os"
    "strings"
)

// 执行评测
func UserRunJudge(c *cli.Context) error {
    if strings.TrimSpace(c.Args().Get(0)) == "" {
        return errors.Errorf("no config file path")
    }
    if strings.TrimSpace(c.Args().Get(1)) == "" {
        return errors.Errorf("no code file path")
    }

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
        client.NewClientSuccessMessage(judgeResult).Print(true)
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
