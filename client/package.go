package client

import (
	"github.com/LanceLRQ/deer-executor/v2/client/packmgr"
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
	&cli.BoolFlag{
		Name:    "zip",
		Aliases: []string{"z"},
		Value:   false,
		Usage:   "Package as a zip file",
	},
}

var AppPackageSubCommands = cli.Commands{
	{
		Name:      "build",
		HelpName:  "deer-executor package build",
		Usage:     "build problem package",
		ArgsUsage: "<configs_file> <output_file>",
		Flags:     PackProblemFlags,
		Action:    packmgr.BuildProblemPackage,
	},
	{
		Name:      "unpack",
		HelpName:  "deer-executor package unpack",
		Usage:     "unpack problem package",
		ArgsUsage: "<package_file> <output_dir>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-validate",
				Value: false,
				Usage: "disable package validation",
			},
		},
		Action: packmgr.UnpackProblemPackage,
	},
	{
		Name:      "info",
		HelpName:  "deer-executor package info",
		Usage:     "show problem package info",
		ArgsUsage: "<package_file>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "gpg",
				Value: false,
				Usage: "output GPG signature info",
			},
		},
		Action: packmgr.ReadProblemInfo,
	},
}
