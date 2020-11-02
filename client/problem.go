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
            &cli.StringFlag{
                Name:  "type",
                Aliases: []string {"t"},
                Value: "all",
                Usage: "module type: validate_cases|test_cases|all",
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
}
