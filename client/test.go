package client

import (
    "fmt"
    "github.com/urfave/cli/v2"
)

func Test(c *cli.Context) error {
    //ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
    //rel, err := utils.RunUnixShell(ctx, "./a.out", []string{""}, func(writer io.Writer) error {
    //    _, err := writer.Write([]byte("1 2\n"))
    //    return err
    //})
    //if err != nil {
    //    return err
    //}
    //fmt.Println(rel.Success)
    //fmt.Println(rel.ExitCode)
    //fmt.Println(rel.Stdout)
    //fmt.Println(rel.Stderr)
    fmt.Println("abc"[:1])
    //fmt.Println("Hello world")
    return nil
}
