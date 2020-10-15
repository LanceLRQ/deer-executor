package main

import (
	"encoding/hex"
	"fmt"
	"github.com/LanceLRQ/deer-executor/client"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/LanceLRQ/deer-executor/persistence"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)


func main() {
	main := &cli.App{
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
				Action: client.Run,
				Flags: client.RunFlags,
			},
			{
				Name: "make",
				Usage: "make configurations",
				Action: func(c *cli.Context) error {
					config := executor.JudgeSession{}
					output := c.String("output")
					if output != "" {
						_, err := os.Stat(output)
						if os.IsExist(err) {
							log.Fatal("output file exists")
							return nil
						}
						fp, err := os.OpenFile(output, os.O_WRONLY | os.O_CREATE, 0644)
						if err != nil {
							log.Fatalf("open output file error: %s\n", err.Error())
							return nil
						}
						defer fp.Close()
						_, err = fp.WriteString(executor.ObjectToJSONStringFormatted(config))
						if err != nil {
							return err
						}
					} else {
						fmt.Println(executor.ObjectToJSONStringFormatted(config))
					}
					return nil
				},
				Flags: []cli.Flag {
					&cli.StringFlag {
						Name: "output",
						Aliases: []string{"out"},
						Usage: "output config file",
					},
				},
			},
			{
				Name: "test",
				Hidden: true,
				Usage: "",
				Action: func(c *cli.Context) error {
					pkey, err := persistence.ReadPemFile("./test/certs/test.key")
					if err != nil { return err }
					sign, err := persistence.RSA2048SignString("Hello World", pkey)
					if err != nil { return err }
					fmt.Println(hex.EncodeToString(sign))

					pukey, err := persistence.ReadPemFile("./test/certs/test.pem")
					if err != nil { return err }
					err = persistence.RSA2048VerifyString("Hello World", sign, pukey)
					if err == nil {
						fmt.Println("Yes!")
					}
					return err
				},
			},
		},
	}
	err := main.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
