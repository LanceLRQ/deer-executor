package main

import (
    "fmt"
    "github.com/LanceLRQ/deer-executor/v2/client"
    "github.com/LanceLRQ/deer-executor/v2/client/generate"
    "github.com/LanceLRQ/deer-executor/v2/client/run"
    "github.com/urfave/cli/v2"
    "os"
)

var (
    buildGitHash   string
    buildTime string
    buildGoVersion string
    buildVersion string
)

func main() {
    main := &cli.App{
        Name:     "Deer Executor",
        HelpName: "deer-executor",
        Usage:    "An executor for online judge.",
        Action: func(c *cli.Context) error {
            fmt.Println("Deer Executor\n--------------------")
            fmt.Printf("version: %s (built %s)\n", buildVersion, buildTime)
            fmt.Printf("commit: %s \n", buildGitHash)
            fmt.Printf("compiler: %s\n", buildGoVersion)
            return nil
        },
        Commands: cli.Commands{
            {
                Name:      "run",
                Usage:     "run code judging",
                Aliases:   []string{"r"},
                ArgsUsage: "<config_file|problem_package> <code_file>",
                Action:    run.UserRunJudge,
                Flags:     client.RunFlags,
            },
            {
                Name:        "example",
                Aliases:     []string{"e"},
                Usage:       "generate all kinds of configuration files",
                Subcommands: client.AppGeneratorSubCommands,
            },
            {
                Name:      "new",
                Aliases:   []string{"n"},
                ArgsUsage: "<output_dir>",
                Usage:     "create a new problem with example",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name:  "name",
                        Value: "a+b",
                        Usage: "template file name (found in ./lib/example/)",
                    },
                },
                Action: generate.InitProblemWorkDir,
            },
            {
                Name:        "package",
                Usage:       "package manager",
                Subcommands: client.AppPackageSubCommands,
            },
            {
                Name:        "problem",
                Aliases:     []string{"p"},
                Usage:       "problem manager",
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
        //client.NewClientErrorMessage(err, nil).Print(true)
        fmt.Printf("%+v\n", err)
        os.Exit(-1)
    }
}
