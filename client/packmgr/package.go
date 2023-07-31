package packmgr

import (
	"encoding/hex"
	"fmt"
	"github.com/LanceLRQ/deer-executor/v3/executor"
	persistence "github.com/LanceLRQ/deer-executor/v3/executor/persistence"
	utils "github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

// BuildProblemPackage build a problem project package with deer-package
func BuildProblemPackage(c *cli.Context) error {

	if c.String("passphrase") != "" {
		log.Println("[warn] Using a password on the command line interface can be insecure.")
	}
	passphrase := []byte(c.String("passphrase"))
	configFile := c.Args().Get(0)
	outputFile := c.Args().Get(1)

	if c.Bool("zip") && !strings.HasSuffix(configFile, "problem.json") {
		return errors.Errorf("config file must named 'problem.json' in zip mode")
	}

	var err error
	var pem *persistence.DigitalSignPEM

	_, err = os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		return errors.Errorf("problem config file (%s) not found", configFile)
	}

	if c.Bool("sign") {
		pem, err = persistence.GetArmorPublicKey(c.String("gpg-key"), passphrase)
		if err != nil {
			return err
		}
	}

	// problem
	session, err := executor.NewSession(configFile)
	if err != nil {
		return err
	}

	pack := persistence.NewProblemProjectPackage(&session.JudgeConfig)
	options := &persistence.ProblemProjectPersisOptions{
		CommonPersisOptions: persistence.CommonPersisOptions{
			DigitalSign: c.Bool("sign"),
			OutFile:     outputFile,
			DigitalPEM:  pem,
		},
		ConfigFile: configFile,
		ProjectDir: session.ConfigDir,
	}

	err = executor.CheckRequireFilesExists(&session.JudgeConfig, options.ProjectDir)
	if err != nil {
		return err
	}
	err = pack.WritePackageFile(options)
	if err != nil {
		return err
	}
	fmt.Println("Done.")
	return nil
}

// UnpackDeerPackage Unpack deer-package
func UnpackDeerPackage(c *cli.Context) error {
	packageFile := c.Args().Get(0)
	workDir := c.Args().Get(1)
	if c.Bool("no-validate") {
		log.Println("[warn] package validation had been disabled!")
	}
	// Check if the file belongs to deer-package
	isDeerPack, packageType, err := utils.IsDeerPackage(packageFile)
	if err != nil {
		return err
	}
	// Unpack deer package
	if !isDeerPack {
		return errors.Errorf("not a deer-executor package file")
	}

	switch packageType {

	case persistence.PackageTypeProblem:
		pack, err := persistence.ParseProblemPackageFile(packageFile, !c.Bool("no-validate"))
		if err != nil {
			return err
		}
		// if <workDir> exists
		if _, err := os.Stat(workDir); err == nil {
			return errors.Errorf("work directory (%s) path exisis", workDir)
		}
		// create folder <workDir>
		if err := os.MkdirAll(workDir, 0775); err != nil {
			return err
		}
		err = pack.UnpackProblemProject(workDir)
		if err != nil {
			return err
		}

	case persistence.PackageTypeJudgeResult:
		pack, err := persistence.ParseJudgeResultPackageFile(packageFile, !c.Bool("no-validate"))
		if err != nil {
			return err
		}
		// if <workDir> exists
		if _, err := os.Stat(workDir); err == nil {
			return errors.Errorf("work directory (%s) path exisis", workDir)
		}
		// create folder <workDir>
		if err := os.MkdirAll(workDir, 0775); err != nil {
			return err
		}
		err = pack.UnpackJudgeResult(workDir)
		if err != nil {
			return err
		}
	}

	fmt.Println("Done.")
	return nil
}

// ReadDeerPackageInfo visit deer-package info
func ReadDeerPackageInfo(c *cli.Context) error {
	packageFile := c.Args().Get(0)
	isDeerPack, packageType, err := utils.IsDeerPackage(packageFile)
	if err != nil {
		return err
	}
	// Check if the file belongs to deer-package
	if !isDeerPack {
		return errors.Errorf("not a deer-executor problem package file")
	}

	printCommonInfo := func(base *persistence.DeerPackageBase) {
		typeName, ok := persistence.PackageTypeNameMap[base.PackageType]
		if !ok {
			typeName = ""
		}
		fmt.Printf(
			"Package ID: %s\nPackage Type: %s\nCommit Version: %d\nPackage File Size: %s\n",
			base.PackageID, typeName, base.CommitVersion, persistence.FormatFileSize(base.GetPackageSize()),
		)
		fmt.Printf("Hash: %s\n", hex.EncodeToString(base.Signature[:]))
		if base.GPGSignSize > 0 {
			fmt.Printf("GPG Signature: Yes (%s...)\n", hex.EncodeToString(base.GPGSignature[:16]))
			g, err := base.GetGPGInfo()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("GPG Info: %s\n", g)
			}
		} else {
			fmt.Printf("GPG Signature: No\n")
		}
	}

	// If it is a DeerPackage file, unpack it.
	switch packageType {

	case persistence.PackageTypeProblem:
		pack, err := persistence.ParseProblemPackageFile(packageFile, !c.Bool("no-validate"))
		if err != nil {
			return err
		}
		printCommonInfo(&pack.DeerPackageBase)
		if c.Bool("json") {
			err = pack.GetProblemConfig()
			if err != nil {
				return err
			}
			fmt.Println(utils.ObjectToJSONStringFormatted(pack.ProblemConfigs))
		}

	case persistence.PackageTypeJudgeResult:
		pack, err := persistence.ParseJudgeResultPackageFile(packageFile, !c.Bool("no-validate"))
		if err != nil {
			return err
		}
		printCommonInfo(&pack.DeerPackageBase)
		if c.Bool("json") {
			err = pack.GetResult()
			if err != nil {
				return err
			}
			fmt.Println(utils.ObjectToJSONStringFormatted(pack.JudgeResult))
		}
	}
	return nil
}
