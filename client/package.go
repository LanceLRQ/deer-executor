package client

import (
	"github.com/LanceLRQ/deer-executor/v3/client/generate"
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
		Aliases:   []string{"b"},
		HelpName:  "deer-executor package build",
		Usage:     "build problem package",
		ArgsUsage: "<configs_file> <output_file>",
		Flags:     PackProblemFlags,
		Action:    packmgr.BuildProblemPackage,
	},
	{
		Name:      "new",
		Aliases:   []string{"n"},
		ArgsUsage: "<output_dir>",
		Usage:     "create a new problem with example",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "sample",
				Value: "",
				Usage: "sample file name or problem package file",
			},
		},
		Action: generate.InitProblemProjectDir,
	},
	{
		Name:      "unpack",
		Aliases:   []string{"u"},
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
		Aliases:   []string{"i"},
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
