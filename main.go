package main

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/client"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	main := &cli.App{
		Name:  "Deer Executor",
		Usage: "An executor for online judge.",
		Action: func(c *cli.Context) error {
			fmt.Println("Deer Executor v2.0")
			return nil
		},
		Commands: cli.Commands{
			{
				Name:      "run",
				Usage:     "run code judging",
				ArgsUsage: "code_file",
				Action:    client.Run,
				Flags:     client.RunFlags,
			},
			{
				Name:  "make",
				Usage: "make/generate somethings",
				Subcommands: cli.Commands{
					{
						Name:   "config",
						Action: client.MakeConfigFile,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "output",
								Aliases: []string{"out"},
								Value: 	 "",
								Usage:   "output config file",
							},
						},
					},
					{
						Name:   "cert",
						Action: client.GenerateRSA,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "bit",
								Value:   2048,
								Aliases: []string{"b"},
								Usage:   "RSA bit",
							},
						},
					},{
						Name:   "compiler",
						Action: client.MakeCompileConfigFile,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "output",
								Aliases: []string{"out"},
								Value: 	 "",
								Usage:   "output config file",
							},
						},
					},{
						Name:   "jit_memory",
						Action: client.MakeJITMemoryConfigFile,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "output",
								Aliases: []string{"out"},
								Value: 	 "",
								Usage:   "output config file",
							},
						},
					},
				},
			},
			{
				Name:   "pack",
				Usage:  "pack problem configs",
				ArgsUsage: "configs_file output_file",
				Flags: client.PackProblemFlags,
				Action: client.PackProblem,
			},
			{
				Name:   "test",
				Hidden: true,
				Usage:  "",
				Action: client.Test,
			},
		},
	}
	err := main.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
