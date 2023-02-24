package server

import "github.com/urfave/cli/v2"

var JudgeServiceCommandFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "./server.yaml",
		Usage:   "Config file",
	},
}
