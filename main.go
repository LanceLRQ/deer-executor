package main

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
)


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
		Name: "special-judge-redirect-program-out",
		Value: true,
		Usage: "Redirect target program's Stdout to checker's Stdin (checker mode). if not, redirect testcase-in file to checker's STDIN",
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
		Name: "session-id",
		Aliases: []string{"s"},
		Value: "",
		Usage: "Custom the session id",
	},
	&cli.StringFlag {
		Name: "work-dir",
		Aliases: []string{"d"},
		Value: "/tmp",
		Usage: "Custom work directory",
	},
	&cli.BoolFlag {
		Name: "clean",
		Aliases: []string{"c"},
		Value: false,
		Usage: "Delete session directory after judge",
	},
}

// create and get session directory
func getSessionDir(workDir string, sessionId string) (string, error) {
	_, err := os.Stat(workDir)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("work dir (%s) not exists", workDir)
	} else if err != nil {
		return "", err
	}
	sessionDir := path.Join(workDir, sessionId)
	_, err = os.Stat(sessionDir)
	if os.IsExist(err) {
		_ = os.RemoveAll(sessionDir)
	}
	err = os.Mkdir(sessionDir, 0755)
	if err != nil {
		return "", fmt.Errorf("create session dir error: %s", err.Error())
	}
	return sessionDir, nil
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
					// create session
					session := executor.JudgeSession{
						CodeFile: c.Args().Get(0),
						CodeLangName: c.String("language"),
						TestCaseIn: c.String("testcase-input"),
						TestCaseOut: c.String("testcase-output"),
						TimeLimit: c.Int("time-limit"),
						MemoryLimit: c.Int("memory-limit"),
						RealTimeLimit: c.Int("real-time-limit"),
						FileSizeLimit: c.Int("file-size-limit"),
						Uid: c.Int("uid"),
						SpecialJudge: executor.SpecialJudgeOptions {
							Mode: c.Int("special-judge"),
							Checker: c.String("special-judge-checker"),
							RedirectProgramOut: c.Bool("special-judge-redirect-program-out"),
							TimeLimit: c.Int("special-judge-time-limit"),
							MemoryLimit: c.Int("special-judge-memory-limit"),
						},
					}
					// fill session id
					if c.String("session") == "" {
						session.SessionId = uuid.NewV1().String()
					} else {
						session.SessionId =  c.String("session")
					}
					if !c.Bool("clean") {
						fmt.Printf("Judge Session: %s\n", session.SessionId)
					}

					// create session dir
					sessionDir, err := getSessionDir(c.String("work-dir"), session.SessionId)
					if err != nil {
						log.Fatal(err)
						return err
					}
					session.SessionDir = sessionDir

					judgeResult, err := session.RunJudge()
					if err != nil {
						log.Fatal(err)
						return err
					}
					fmt.Println(executor.ObjectToJSONStringFormatted(judgeResult))

					// Do clean
					if c.Bool("clean") {
						_ = os.RemoveAll(sessionDir)
					}
					return err
				},
				Flags: commonFlags,
			},
		},
	}).Run(os.Args)
}
