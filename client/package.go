package client

import (
	"crypto/rsa"
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/LanceLRQ/deer-executor/persistence"
	"github.com/LanceLRQ/deer-executor/persistence/problems"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/openpgp"
	"log"
	"os"
)

var PackProblemFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:     	"sign",
		Aliases:  	[]string{"s"},
		Value: 	 	false,
		Usage:    	"Enable digital sign (GPG)",
	},
	&cli.StringFlag {
		Name: 		"private-key",
		Aliases:  	[]string{"pri"},
		Value: 		"",
		Usage: 		"Digital sign private key",
	},
	&cli.StringFlag {
		Name: 		"public-key",
		Aliases:  	[]string{"pub"},
		Value: 		"",
		Usage: 		"Digital sign public key",
	},
	&cli.StringFlag {
		Name: 		"passphrase",
		Aliases:  	[]string{"p", "password", "pwd"},
		Value: 		"",
		Usage: 		"Private key passphrase",
	},
}

func PackProblem(c *cli.Context) error {

	if c.String("passphrase") != "" {
		log.Println("[warn] Using a password on the command line interface can be insecure.")
	}
	passphrase := []byte(c.String("passphrase"))
	configFile := c.Args().Get(0)
	outputFile := c.Args().Get(1)

	// Open key
	keyRingReader, err := os.Open(c.String("private-key"))
	if err != nil {
		return err
	}
	// Read GPG Keys
	elist, err := openpgp.ReadArmoredKeyRing(keyRingReader)
	if err != nil {
		return err
	}
	if len(elist) < 1 {
		return fmt.Errorf("file has no GPG key")
	}
	gpgKey := elist[0].PrivateKey
	if gpgKey.Encrypted {
		if len(passphrase) == 0 {
			passphrase, err = gopass.GetPasswdPrompt("please input passphrase of key> ", true, os.Stdin, os.Stdout)
			if err != nil {
				return err
			}
		}
		err = gpgKey.Decrypt(passphrase)
		if err != nil {
			return err
		}
	}

	// problem
	session, err := executor.NewSession(configFile)
	if err != nil {
		return err
	}
	pem := persistence.DigitalSignPEM {
		PrivateKey: gpgKey.PrivateKey.(*rsa.PrivateKey),
		PublicKeyRaw: []byte{},
		PublicKey: elist[0].PrimaryKey.PublicKey.(*rsa.PublicKey),
	}
	options := problems.ProblemPersisOptions{
		DigitalSign: true,
		DigitalPEM: pem,
		OutFile: outputFile,
	}
	err = problems.PackProblems(session, options)
	if err != nil {
		return err
	}

	return nil
}
