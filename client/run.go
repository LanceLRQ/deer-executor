package client

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
	"time"
)

var RunFlags = []cli.Flag {
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
		Aliases: []string{"rtl"},
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
		Usage: "Executable checker file or checker's source code",
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
		Name: "session-root",
		Aliases: []string{"sr"},
		Value: "/tmp",
		Usage: "Custom work directory",
	},
	&cli.BoolFlag {
		Name: "clean",
		Aliases: []string{"c"},
		Value: false,
		Usage: "Delete session directory after judge",
	},
	&cli.BoolFlag {
		Name: "debug",
		Value: false,
		Usage: "print debug log",
	},
	&cli.IntFlag {
		Name: "benchmark",
		Value: 0,
		Usage: "Start benchmark",
	},
}

func run(c *cli.Context, counter int) (*executor.JudgeResult, error) {
	isBenchmarkMode := c.Int("benchmark") > 1

	// create session
	session := executor.JudgeSession{
		CodeFile: c.Args().Get(0),
		CodeLangName: c.String("language"),
		TestCases: []executor.TestCase {
			{
				Id: "test",
				TestCaseIn: c.String("testcase-input"),
				TestCaseOut: c.String("testcase-output"),
			},
		},
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
	// Do clean (or benchmark on)
	if c.Bool("clean") || isBenchmarkMode {
		defer session.Clean()
	}
	// fill session id
	if isBenchmarkMode {
		session.SessionId = uuid.NewV1().String() + strconv.Itoa(counter)
	} else {
		if c.String("session") == "" {
			session.SessionId = uuid.NewV1().String()
		} else {
			session.SessionId = c.String("session")
		}
	}

	if !c.Bool("clean") && c.Int("benchmark") <= 1 {
		log.Println(fmt.Sprintf("Judge Session: %s\n", session.SessionId))
	}

	// create session dir
	session.SessionRoot = c.String("session-root")
	sessionDir, err := getSessionDir(session.SessionRoot, session.SessionId)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	session.SessionDir = sessionDir

	judgeResult := session.RunJudge()

	return &judgeResult, err
}

func Run(c *cli.Context) error {
	n := c.Int("benchmark")
	if n <= 1 {
		judgeResult, err := run(c, 0)
		if err != nil {
			return err
		}
		fmt.Println(executor.ObjectToJSONStringFormatted(judgeResult))
	} else {
		rfp, err := os.OpenFile("./report.log", os.O_WRONLY | os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer rfp.Close()

		startTime := time.Now().UnixNano()
		exitCounter := map[int]int {}
		for i := 0; i < n; i++ {
			if i % 10 == 0 {
				fmt.Printf("%d / %d\n", i, n)
			}
			judgeResult, err := run(c, i)
			if err != nil {
				fmt.Printf("break! %s\n", err.Error())
				_, _ = rfp.WriteString(fmt.Sprintf("[%s]: %s\n", strconv.Itoa(i), err.Error()))
				break
			}
			name, ok := executor.FlagMeansMap[judgeResult.JudgeResult]
			if !ok { name = "Unknown" }
			_, _ = rfp.WriteString(fmt.Sprintf("[%s]: %s\n", judgeResult.SessionId, name))
			if judgeResult.JudgeResult != executor.JudgeFlagAC {
				_, _ = rfp.WriteString(executor.ObjectToJSONStringFormatted(judgeResult) + "\n")
			}
			exitCounter[judgeResult.JudgeResult]++
		}
		endTime := time.Now().UnixNano()
		for key, value := range exitCounter {
			name, ok := executor.FlagMeansMap[key]
			if !ok { name = "Unknown" }
			fmt.Printf("%s: %d\n", name, value)
		}
		duration := float64(endTime - startTime) / float64(time.Second)
		fmt.Printf("total time used: %.2fs\n", duration)
	}
	return nil
}