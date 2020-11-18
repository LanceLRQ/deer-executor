package run

import (
    "fmt"
    "github.com/urfave/cli/v2"
)


// 执行评测
func UserRunJudge(c *cli.Context) error {
    fmt.Println("Sorry, runner only supporting linux/darwin")
    return nil
}
