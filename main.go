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
				Name:   "test",
				Hidden: true,
				Usage:  "",
				Action: func(c *cli.Context) error {
					//privateKey, err := persistence.ReadPemFile("./data/certs/test.key")
					//if err != nil {
					//	return err
					//}
					//sign, err := persistence.RSA2048SignString("Hello World", pkey)
					//if err != nil {
					//	return err
					//}
					//fmt.Println(hex.EncodeToString(sign))
					//
					//publicKey, err := persistence.ReadPemFile("./data/certs/test.pem")
					//if err != nil {
					//	return err
					//}
					//err = persistence.RSA2048VerifyString("Hello World", sign, publicKey)
					//if err == nil {
					//	fmt.Println("Yes!")
					//}

					//rel, err := persistence.SHA256Streams([]io.Reader{
					//	bytes.NewReader(publicKey),
					//	bytes.NewReader(privateKey),
					//})
					//if err != nil {
					//	return err
					//}
					//fmt.Println(hex.EncodeToString(rel))
					//_, err := judge_result.ReadJudgeResult("./result")
					//if err != nil {
					//	return err
					//}
					//fmt.Println(executor.ObjectToJSONStringFormatted(rst))
					return nil
				},
			},
		},
	}
	err := main.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
