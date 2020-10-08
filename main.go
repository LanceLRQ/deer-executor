package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	(&cli.App{
		Name: "Deer Executor",
		Usage: "An executor for online judge.",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello World")
			return nil
		},
		Commands: cli.Commands {
			{
				Name: "run",
				Usage: "run code judging",
				ArgsUsage: "code_file",
				Action: func(c *cli.Context) error {
					fmt.Println("Test")
					return nil
				},
				Flags: []cli.Flag {
					&cli.StringFlag {
						Name: "testcase-input",
						Required: true,
						Aliases: []string{"tin"},
						Usage: "Testcase input file",
					},
					&cli.StringFlag {
						Name: "testcase-output",
						Aliases: []string{"tout"},
						Required: true,
						Usage: "Testcase output file",
					},
					&cli.StringFlag {
						Name: "program-output",
						Value: "/tmp/program.out",
						Aliases: []string{"pout"},
						Usage: "Program stdout file",
					},
					&cli.StringFlag {
						Name: "program-stderr",
						Value: "/tmp/program.err",
						Aliases: []string{"perr"},
						Usage: "Program stderr file",
					},
					&cli.IntFlag {
						Name: "time-limit",
						Value: 1000,
						Aliases: []string{"tl"},
						Usage: "Time limit (ms)",
					},
					&cli.IntFlag {
						Name: "memory-limit",
						Value: 65535,
						Aliases: []string{"ml"},
						Usage: "Memory limit (KB)",
					},
					&cli.IntFlag {
						Name: "real-time-limit",
						Value: 0,
						Usage: "Real Time Limit (ms)",
					},
					&cli.IntFlag {
						Name: "file-size-limit",
						Value: 100 * 1024 * 1024,
						Usage: "File Size Limit (bytes)",
					},
					&cli.IntFlag {
						Name: "uid",
						Value: -1,
						Usage: "User id",
					},
					&cli.StringFlag {
						Name: "language",
						Aliases: []string{"lang"},
						Value: "auto",
						Usage: "Coding language",
					},
				},
			},
		},
	}).Run(os.Args)
}

func RunJudge() {
	//compiler := provider.GnucCompileProvider{}
}