package client

import (
    "github.com/LanceLRQ/deer-executor/client/problem"
    "github.com/urfave/cli/v2"
)

var PackProblemFlags = []cli.Flag{
    &cli.BoolFlag{
        Name:    "sign",
        Aliases: []string{"s"},
        Value:   false,
        Usage:   "Enable digital sign (GPG)",
    },
    &cli.StringFlag{
        Name:    "gpg-key",
        Aliases: []string{"key"},
        Value:   "",
        Usage:   "GPG private key file",
    },
    &cli.StringFlag{
        Name:    "passphrase",
        Aliases: []string{"p", "password", "pwd"},
        Value:   "",
        Usage:   "GPG private key passphrase",
    },
}

var AppPackageSubCommands = cli.Commands{
    {
        Name:      "build",
        Usage:     "build problem package",
        ArgsUsage: "configs_file output_file",
        Flags:     PackProblemFlags,
        Action:    problem.BuildProblemPackage,
    },
    {
        Name:      "info",
        Usage:     "show problem package info",
        ArgsUsage: "package_file",
        Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:  "sign",
                Value: false,
                Usage: "output GPG signature info",
            },
        },
        Action: problem.ReadProblemInfo,
    },
}

