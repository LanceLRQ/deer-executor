package client

import (
	"github.com/LanceLRQ/deer-executor/v3/client/packmgr"
	"github.com/urfave/cli/v2"
)

// PackProblemFlags for cli command 'pack'
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

// AppPackageSubCommands for cli command 'pack'
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
		Usage:     "unpack problem / judge_result package",
		ArgsUsage: "<package_file> <output_dir>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-validate",
				Value: false,
				Usage: "disable package validation",
			},
		},
		Action: packmgr.UnpackDeerPackage,
	},
	{
		Name:      "info",
		HelpName:  "deer-executor package info",
		Usage:     "show deer package info",
		ArgsUsage: "<package_file>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Value: false,
				Usage: "output json info (problem config or judge result)",
			},
		},
		Action: packmgr.ReadDeerPackageInfo,
	},
}
