package client

import (
    "github.com/LanceLRQ/deer-executor/packmgr"
    "github.com/urfave/cli/v2"
)

var AppProblemSubCommands = cli.Commands{
    {
        Name:      "build",
        Aliases:   []string{"b"},
        Usage:     "compile binary source codes",
        ArgsUsage: "configs_file",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:  "library",
                Aliases: []string {"l"},
                Value: "./lib",
                Usage: "library dir, contains \"testlib.h\" and \"bits/stdc++.h\" etc.",
            },
        },
        Action: packmgr.CompileProblemWorkDirSourceCodes,
    },
    {
        Name:      "validate",
        Aliases:   []string{"v"},
        Usage:     "validate input case",
        ArgsUsage: "configs_file",
        Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:  "silence",
                Aliases: []string {"s"},
                Value: false,
                Usage: "silence mode",
            },
            &cli.StringFlag{
                Name:  "type",
                Aliases: []string {"t"},
                Value: "all",
                Usage: "module type: validator_cases|test_cases|all",
            },
            &cli.IntFlag{
                Name:  "case",
                Aliases: []string {"c"},
                Value: -1,
                Usage: "case index, -1 means all. when module type set 'all'，it would't work.",
            },
        },
        Action: packmgr.RunTestlibValidators,
    },
    {
        Name:      "generate",
        Aliases:   []string{"gen", "g"},
        Usage:     "generate test case's input/output",
        ArgsUsage: "configs_file",
        Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:  "silence",
                Aliases: []string {"s"},
                Value: false,
                Usage: "silence mode",
            },
            &cli.BoolFlag{
                Name:  "with-answer",
                Usage: "generate answer",
            },
            &cli.UintFlag{
                Name:  "answer",
                Value: 0,
                Usage: "answer case index.",
            },
            &cli.IntFlag{
                Name:  "case",
                Aliases: []string {"c"},
                Value: -1,
                Usage: "case index, -1 means all. when module type set 'all'，it would't work.",
            },
        },
        Action: packmgr.RunTestCaseGenerator,
    },
}
