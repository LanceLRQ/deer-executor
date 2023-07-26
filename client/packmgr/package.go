package packmgr

import (
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

// BuildProblemPackage  创建题目数据包
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

// UnpackProblemPackage 题目包解包
func UnpackProblemPackage(c *cli.Context) error {
	packageFile := c.Args().Get(0)
	workDir := c.Args().Get(1)
	if c.Bool("no-validate") {
		log.Println("[warn] package validation had been disabled!")
	}
	// 检查是否为题目包
	isDeerPack, err := utils.IsDeerPackage(packageFile)
	if err != nil {
		return err
	}
	// 解包
	if isDeerPack {
		pack, err := persistence.ParsePackageFile(packageFile, !c.Bool("no-validate"))
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
	} else {
		return errors.Errorf("not a deer-executor problem package file")
	}
	fmt.Println("Done.")
	return nil
}

// ReadProblemInfo 访问题目包信息
func ReadProblemInfo(c *cli.Context) error {
	packageFile := c.Args().Get(0)
	isDeerPack, err := utils.IsDeerPackage(packageFile)
	if err != nil {
		return err
	}
	// 如果是题目包文件，进行解包
	if isDeerPack {
		pack, err := persistence.ParsePackageFile(packageFile, !c.Bool("no-validate"))
		if err != nil {
			return err
		}
		err = pack.GetProblemConfig()
		if err != nil {
			return err
		}
		if c.Bool("gpg") {
			g, err := pack.GetProblemGPGInfo()
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
			fmt.Println(g)
		} else {
			fmt.Println(utils.ObjectToJSONStringFormatted(pack.ProblemConfigs))
		}
	} else {
		return errors.Errorf("not a deer-executor problem package file")
	}
	return nil
}
