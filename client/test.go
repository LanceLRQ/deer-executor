package client

import (
    "fmt"
    "github.com/urfave/cli/v2"
)

func Test(c *cli.Context) error {
    fmt.Println("Hello world")
    return nil
}
