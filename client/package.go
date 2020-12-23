package client

import (
    "github.com/LanceLRQ/deer-executor/client/packmgr"
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
        Action:    packmgr.BuildProblemPackage,
    },
    {
        Name:      "unpack",
        Usage:     "unpack problem package",
        ArgsUsage: "package_file output_dir",
        Flags:     []cli.Flag{
            &cli.BoolFlag{
                Name:  "no-validate",
                Value: false,
                Usage: "disable package validation",
            },
        },
        Action:    packmgr.UnpackProblemPackage,
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
        Action: packmgr.ReadProblemInfo,
    },
}

