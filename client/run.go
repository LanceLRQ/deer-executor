package client

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/LanceLRQ/deer-executor/persistence"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
	"time"
)

var RunFlags = []cli.Flag {
	&cli.StringFlag {
		Name: "config",
		Aliases: []string{"c"},
		Required: true,
		Usage: "Load configuration file",
	},
	&cli.BoolFlag {
		Name: "clean",
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
	&cli.StringFlag {
		Name: "persistence",
		Aliases: []string{"p"},
		Value: "",
		Usage: "Persistent judge result to file (support: gzip, none)",
	},
	&cli.StringFlag {
		Name: "compressor",
		Value: "gzip",
		Usage: "Persistent compressor type",
	},
	&cli.BoolFlag {
		Name: "digital-sign",
		Aliases: []string{"sign"},
		Value: false,
		Usage: "Enable digital sign",
	},
	&cli.StringFlag {
		Name: "public-key",
		Value: "",
		Usage: "Digital sign public key",
	},
	&cli.StringFlag {
		Name: "private-key",
		Value: "",
		Usage: "Digital sign private key",
	},
}

func run(c *cli.Context, counter int) (*executor.JudgeResult, error) {
	isBenchmarkMode := c.Int("benchmark") > 1
	persistenceOn := c.String("persistence") != ""
	digitalSign := c.Bool("digital-sign")
	compressorType := uint8(1)
	if c.String("compressor") == "none" {
		compressorType = uint8(0)
	}
	jOption := persistence.JudgeResultPersisOptions{
		OutFile: c.String("persistence"),
		CompressorType: compressorType,
		DigitalSign: digitalSign,
	}
	if !isBenchmarkMode && persistenceOn {
		if digitalSign {
			if c.String("public-key") == "" || c.String("private-key") == "" {
				return nil, fmt.Errorf("digital sign need public key and private key")
			}
			digPEM, err := persistence.GetDigitalPEMFromFile(c.String("public-key"), c.String("private-key"))
			if err != nil {
				return nil, err
			}
			if digPEM.PrivateKey == nil || digPEM.PublicKey == nil {
				return nil, fmt.Errorf("parse public key or private key error")
			}
			jOption.DigitalPEM = *digPEM
		}
	}
	// create session
	session := executor.NewSession()
	// laod configuration
	if c.String("config") != "" {
		cbody, err := executor.ReadFile(c.String("config"))
		if err != nil {
			return nil, err
		}
		executor.JSONBytesObject(cbody, session)
	}
	// init files
	session.CodeFile = c.Args().Get(0)
	// Do clean (or benchmark on)
	if c.Bool("clean") || isBenchmarkMode {
		defer session.Clean()
	}
	// create session info
	if isBenchmarkMode {
		session.SessionId = uuid.NewV1().String() + strconv.Itoa(counter)
	} else {
		if session.SessionId == "" {
			session.SessionId = uuid.NewV1().String()
		}
	}
	if session.SessionRoot != "" {
		session.SessionRoot = "/tmp"
	}
	//if !c.Bool("clean") && c.Int("benchmark") <= 1 {
	//	log.Println(fmt.Sprintf("Judge Session: %s\n", session.SessionId))
	//}
	sessionDir, err := getSessionDir(session.SessionRoot, session.SessionId)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	session.SessionDir = sessionDir
	// start judge
	judgeResult := session.RunJudge()
	// persistence
	if !isBenchmarkMode && persistenceOn {
		err = persistence.PersistentJudgeResult(&judgeResult, jOption)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	return &judgeResult, nil
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