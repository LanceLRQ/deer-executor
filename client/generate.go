package client

import (
	"github.com/LanceLRQ/deer-executor/v2/client/generate"
	"github.com/urfave/cli/v2"
)

var AppGeneratorSubCommands = cli.Commands{
	{
		Name:     "config",
		HelpName: "deer-executor example config",
		Action:   generate.MakeProblemConfigFile,
		Usage:    "generate problem config file",
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
		Name:     "compiler",
		HelpName: "deer-executor example compiler",
		Action:   generate.MakeCompileConfigFile,
		Usage:    "generate compiler settings file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"out"},
				Value:   "",
				Usage:   "output config file",
			},
		},
	}, {
		Name:     "jit_memory",
		HelpName: "deer-executor example jit_memory",
		Action:   generate.MakeJITMemoryConfigFile,
		Usage:    "generate jit memory limitation settings file",
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
