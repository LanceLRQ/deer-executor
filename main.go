package main

import (
    "fmt"
    "github.com/LanceLRQ/deer-executor/client"
    "github.com/LanceLRQ/deer-executor/client/generate"
    "github.com/LanceLRQ/deer-executor/client/run"
    "github.com/urfave/cli/v2"
    "log"
    "os"
)

func main() {
    main := &cli.App{
        Name:  "Deer Executor",
        Usage: "An executor for online judge.",
        Action: func(c *cli.Context) error {
            fmt.Println("Deer Executor v2.0")
            return nil
        },
        Commands: cli.Commands{
            {
                Name:      "run",
                Usage:     "run code judging",
                Aliases:   []string{"r"},
                ArgsUsage: "<config_file/problem_package> <code_file>",
                Action:    run.UserRunJudge,
                Flags:     client.RunFlags,
            },
            {
                Name:        "generate",
                Aliases:     []string{"g"},
                Usage:       "generate problem config, compiler config or jit-memory config file",
                Subcommands: client.AppGeneratorSubCommands,
            },
            {
                Name:      "init",
                Aliases:   []string{"init"},
                ArgsUsage: "<config_file>",
                Usage:     "init problem work directory",
                Flags:  []cli.Flag {
                    &cli.StringFlag{
                        Name: "example",
                        Value: "",
                        Usage: "template file (found in ./lib/example/)",
                    },
                },
                Action:    generate.InitProblemWorkDir,
            },
            {
                Name:        "package",
                Aliases:     []string{"a"},
                Usage:       "problem package manager",
                Subcommands: client.AppPackageSubCommands,
            },
            {
                Name:        "problem",
                Aliases:     []string{"p"},
                Usage:       "problem workdir manager",
                Subcommands: client.AppProblemSubCommands,
            },
            {
                Name:   "test",
                Hidden: true,
                Usage:  "",
                Action: client.Test,
            },
        },
    }
    err := main.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
