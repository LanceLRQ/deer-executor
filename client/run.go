package client

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/LanceLRQ/deer-executor/persistence"
	"github.com/LanceLRQ/deer-executor/persistence/judge_result"
	"github.com/LanceLRQ/deer-executor/persistence/problems"
	"github.com/LanceLRQ/deer-executor/provider"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
	"time"
)

var RunFlags = []cli.Flag {
	&cli.BoolFlag {
		Name: "clean",
		Value: false,
		Usage: "Delete session directory after judge",
	},
	&cli.StringFlag {
		Name: "language",
		Aliases: []string{"l"},
		Value: "auto",
		Usage: "Code language name",
	},
	&cli.BoolFlag {
		Name: "debug",
		Value: false,
		Usage: "Print debug log",
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
		Name: "compress",
		Value: "gzip",
		Usage: "Persistent compressor type",
	},
	&cli.BoolFlag {
		Name: "digital-sign",
		Aliases: []string{"sign"},
		Value: false,
		Usage: "Enable digital sign",
	},
	&cli.BoolFlag {
		Name: "detail",
		Value: false,
		Usage: "Show test-cases details",
	},
	&cli.StringFlag {
		Name: "public-key",
		Aliases: []string{"pub"},
		Value: "",
		Usage: "Digital sign public key",
	},
	&cli.StringFlag {
		Name: "private-key",
		Aliases: []string{"pri"},
		Value: "",
		Usage: "Digital sign private key",
	},
	&cli.StringFlag {
		Name: "work-dir",
		Aliases:  []string{"w"},
		Value: "",
		Usage: "Working dir, using to unpack problem package",
	},
}

func run(c *cli.Context, configFile string, counter int) (*executor.JudgeResult, *executor.JudgeSession, error) {
	isBenchmarkMode := c.Int("benchmark") > 1
	// create session
	session, err := executor.NewSession(configFile)
	if err != nil {
		return nil, nil, err
	}
	if c.String("language") != "" {
		session.CodeLangName = c.String("language")
	}
	// init files
	session.CodeFile = c.Args().Get(1)
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
	sessionDir, err := GetSessionDir(session.SessionRoot, session.SessionId)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	session.SessionDir = sessionDir
	// start judge
	judgeResult := session.RunJudge()
	return &judgeResult, session, nil
}

func Run(c *cli.Context) error {
	// 载入默认配置
	err := provider.PlaceCompilerCommands("./compilers.json")
	if err != nil {
		return err
	}
	err = executor.PlaceMemorySizeForJIT("./jit_memory.json")
	if err != nil {
		return err
	}

	configFile := c.Args().Get(0)
	yes, err := problems.IsProblemPackage(configFile)
	if err != nil {
		return err
	}
	// 如果是题目包文件，进行解包
	if yes {
		workDir := c.String("work-dir")
		autoRemoveWorkDir := false;
		if workDir == "" {
			workDir = "/tmp/" + uuid.NewV4().String()
			autoRemoveWorkDir = true
		}
		if info, err := os.Stat(workDir); os.IsNotExist(err) {
			err = os.MkdirAll(workDir, 0755)
			if err != nil {
				return err
			}
		} else if !info.IsDir() {
			return fmt.Errorf("work dir path cannot be a file path")
		}
		s, err := problems.ReadProblemInfo(configFile, true, workDir)
		if err != nil {
			return err
		}
		configFile = s.ConfigFile
		if autoRemoveWorkDir {
			defer (func() {
				_ = os.RemoveAll(workDir)
			})()
		}
	}
	isBenchmarkMode := c.Int("benchmark") > 1
	benchmarkN := c.Int("benchmark")
	if !isBenchmarkMode {
		// 正常运行
		// parse params
		persistenceOn := c.String("persistence") != ""
		digitalSign := c.Bool("digital-sign")
		compressorType := uint8(1)
		if c.String("compress") == "none" {
			compressorType = uint8(0)
		}
		jOption := judge_result.JudgeResultPersisOptions{
			OutFile: c.String("persistence"),
			CompressorType: compressorType,
			DigitalSign: digitalSign,
		}
		// 是否要持久化结果
		if persistenceOn {
			if digitalSign {
				if c.String("public-key") == "" || c.String("private-key") == "" {
					return fmt.Errorf("digital sign need public key and private key")
				}
				digPEM, err := persistence.GetDigitalPEMFromFile(c.String("public-key"), c.String("private-key"))
				if err != nil {
					return err
				}
				if digPEM.PrivateKey == nil || digPEM.PublicKey == nil {
					return fmt.Errorf("parse public key or private key error")
				}
				jOption.DigitalPEM = *digPEM
			}
		}
		// Start Judge
		judgeResult, judgeSession, err := run(c, configFile, 0)
		if err != nil {
			return err
		}
		// Do clean (or benchmark on)
		if c.Bool("clean") || isBenchmarkMode {
			defer judgeSession.Clean()
		}
		// persistence
		if !isBenchmarkMode && persistenceOn {
			err = judge_result.PersistentJudgeResult(judgeSession, judgeResult, jOption)
			if err != nil {
				return err
			}
		}
		if !c.Bool("detail") {
			judgeResult.TestCases = nil
		}
		fmt.Println(executor.ObjectToJSONStringFormatted(judgeResult))
		os.Exit(judgeResult.JudgeResult)
	} else {
		// 基准测试
		rfp, err := os.OpenFile("./report.log", os.O_WRONLY | os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer rfp.Close()

		startTime := time.Now().UnixNano()
		exitCounter := map[int]int {}
		for i := 0; i < benchmarkN; i++ {
			if i % 10 == 0 {
				fmt.Printf("%d / %d\n", i, benchmarkN)
			}
			judgeResult, _, err := run(c, configFile, i)
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