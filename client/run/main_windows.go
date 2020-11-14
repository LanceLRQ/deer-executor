package run

import (
    "fmt"
    "github.com/urfave/cli/v2"
)


// 执行评测
func UserRunJudge(c *cli.Context) error {
    fmt.Println("The functional is base on linux/darwin")
    return nil
}
