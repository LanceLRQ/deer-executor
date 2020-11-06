package client

import (
    "github.com/LanceLRQ/deer-executor/client/generate"
    "github.com/urfave/cli/v2"
)

var AppMakeSubCommands = cli.Commands{
    {
        Name:   "config",
        Action: generate.MakeProblemConfigFile,
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:    "output",
                Aliases: []string{"out"},
                Value:   "",
                Usage:   "output config file",
            },
        },
    },
    {
        Name:   "compiler",
        Action: generate.MakeCompileConfigFile,
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:    "output",
                Aliases: []string{"out"},
                Value:   "",
                Usage:   "output config file",
            },
        },
    }, {
        Name:   "jit_memory",
        Action: generate.MakeJITMemoryConfigFile,
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:    "output",
                Aliases: []string{"out"},
                Value:   "",
                Usage:   "output config file",
            },
        },
    },
}
