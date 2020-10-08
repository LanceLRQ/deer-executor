package main

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/docker/docker/pkg/reexec"
	"github.com/urfave/cli/v2"
	"os"
)

func init() {
	reexec.Register("targetProgram", executor.RunTargetProgramProcess)
	reexec.Register("judgeProgram", executor.RunTargetProgramProcess)
	if reexec.Init() {
		os.Exit(0)
	}
}


func RunJudge(options executor.JudgeOptions) error {
	//compiler := provider.GnucCompileProvider{}
	return nil
}


var commonFlags = []cli.Flag {
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
	&cli.IntFlag{
		Name: "special-judge",
		Aliases: []string{"mode"},
		Value: 0,
		Usage: "Special Judge Mode: 0-Disabled；1-Normal；2-Interactor",
	},
	&cli.StringFlag {
		Name: "special-judge-checker",
		Aliases: []string{"checker"},
		Value: "",
		Usage: "Checker file path",
	},
	&cli.BoolFlag {
		Name: "special-judge-redirect-std",
		Value: true,
		Usage: "Redirect target program's Stdout to checker's Stdin (checker mode)",
	},
	&cli.IntFlag {
		Name: "special-judge-time-limit",
		Aliases: []string{"spj-tl"},
		Value: 1000,
		Usage: "Special judge Time limit (ms)",
	},
	&cli.IntFlag {
		Name: "special-judge-memory-limit",
		Aliases: []string{"spj-ml"},
		Value: 65535,
		Usage: "Special judge memory limit (kb)",
	},
	&cli.StringFlag {
		Name: "special-judge-checker-stdout",
		Aliases: []string{"cout"},
		Value: "/tmp/checker.out",
		Usage: "Special judge checker's stdout",
	},
	&cli.StringFlag {
		Name: "special-judge-checker-stderr",
		Aliases: []string{"cerr"},
		Value: "/tmp/checker.err",
		Usage: "Special judge checker's stderr",
	},
	&cli.StringFlag {
		Name: "special-judge-checker-logfile",
		Aliases: []string{"log"},
		Value: "/tmp/judge.log",
		Usage: "Special judge checker's log file params",
	},
}

func main() {
	(&cli.App{
		Name: "Deer Executor",
		Usage: "An executor for online judge.",
		Action: func(c *cli.Context) error {
			fmt.Println("Deer Executor v2.0")
			return nil
		},
		Commands: cli.Commands {
			{
				Name: "run",
				Usage: "run code judging",
				ArgsUsage: "code_file",
				Action: func(c *cli.Context) error {
					options := executor.JudgeOptions{
						CodeFile: c.Args().Get(0),
						CodeLangName: c.String("language"),
						TestCaseIn: c.String("testcase-input"),
						TestCaseOut: c.String("testcase-output"),
						ProgramOut: c.String("program-output"),
						ProgramError: c.String("program-stderr"),
						TimeLimit: c.Int("time-limit"),
						MemoryLimit: c.Int("memory-limit"),
						RealTimeLimit: c.Int("real-time-limit"),
						FileSizeLimit: c.Int("file-size-limit"),
						Uid: c.Int("uid"),
						SpecialJudge: executor.SpecialJudgeOptions {
							Mode: c.Int("special-judge"),
							Checker: c.String("special-judge-checker"),
							RedirectStd: c.Bool("special-judge-redirect-std"),
							TimeLimit: c.Int("special-judge-time-limit"),
							MemoryLimit: c.Int("special-judge-memory-limit"),
							Stdout: c.String("special-judge-checker-stdout"),
							Stderr: c.String("special-judge-checker-stderr"),
							Logfile: c.String("special-judge-checker-logfile"),
						},
					}
					return RunJudge(options)
				},
				Flags: commonFlags,
			},
		},
	}).Run(os.Args)
}
