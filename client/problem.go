package client

import (
    "github.com/LanceLRQ/deer-executor/packmgr"
    "github.com/urfave/cli/v2"
)

var AppProblemSubCommands = cli.Commands{
    {
        Name:      "build",
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
}

