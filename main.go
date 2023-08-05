package main

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/v3/agent"
	agentConfig "github.com/LanceLRQ/deer-executor/v3/agent/config"
	"github.com/LanceLRQ/deer-executor/v3/client"
	"github.com/LanceLRQ/deer-executor/v3/client/run"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	buildGitHash   string
	buildTime      string
	buildGoVersion string
	buildVersion   string
)

func main() {
	main := &cli.App{
		Name:     "Deer Executor",
		HelpName: "deer-executor",
		Usage:    "An executor for online judge.",
		Writer:   os.Stderr,
		Action: func(c *cli.Context) error {
			fmt.Println("Deer Executor\n--------------------")
			fmt.Printf("version: %s (built %s)\n", buildVersion, buildTime)
			fmt.Printf("commit: %s \n", buildGitHash)
			fmt.Printf("compiler: %s\n", buildGoVersion)
			return nil
		},
		Commands: cli.Commands{
			{
				Name:      "run",
				Usage:     "run code judging",
				Aliases:   []string{"r"},
				ArgsUsage: "<config_file|problem_package> <code_file>",
				Action:    run.UserRunJudge,
				Flags:     client.RunFlags,
			},
			{
				Name:        "example",
				Aliases:     []string{"e"},
				Usage:       "generate all kinds of configuration files",
				Subcommands: client.AppGeneratorSubCommands,
			},
			{
				Name:        "package",
				Usage:       "package manager",
				Subcommands: client.AppPackageSubCommands,
			},
			{
				Name:        "problem",
				Aliases:     []string{"p"},
				Usage:       "problem manager",
				Subcommands: client.AppProblemSubCommands,
			},
			{
				Name:   "agent",
				Hidden: true,
				Usage:  "judge service agent service",
				Before: func(c *cli.Context) error {
					config.AddDriver(yamlv3.Driver)
					err := config.LoadFiles(c.String("config"))
					if err != nil {
						return err
					}
					// 载入服务端配置
					return agentConfig.LoadGlobalConf()
				},
				Action: agent.LaunchJudgeService,
				Flags:  agent.JudgeServiceCommandFlags,
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
