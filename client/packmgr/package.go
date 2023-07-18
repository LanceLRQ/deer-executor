package packmgr

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/v3/executor"
	persistence2 "github.com/LanceLRQ/deer-executor/v3/executor/persistence"
	problems2 "github.com/LanceLRQ/deer-executor/v3/executor/persistence/problems"
	utils2 "github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
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
	var pem *persistence2.DigitalSignPEM

	_, err = os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		return errors.Errorf("problem config file (%s) not found", configFile)
	}

	if c.Bool("sign") {
		pem, err = persistence2.GetArmorPublicKey(c.String("gpg-key"), passphrase)
		if err != nil {
			return err
		}
	}
	options := persistence2.ProblemPackageOptions{}
	options.ConfigFile = configFile
	options.DigitalSign = c.Bool("sign")
	options.DigitalPEM = pem
	options.OutFile = outputFile

	// problem
	session, err := executor.NewSession(configFile)
	if err != nil {
		return err
	}
	options.ConfigDir = session.ConfigDir

	err = executor.CheckRequireFilesExists(&session.JudgeConfig, options.ConfigDir)
	if err != nil {
		return err
	}

	if c.Bool("zip") {
		err = problems2.PackProblemsAsZip(&options)
		if err != nil {
			return err
		}
	} else {
		err = problems2.PackProblems(&session.JudgeConfig, &options)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnpackProblemPackage 题目包解包
func UnpackProblemPackage(c *cli.Context) error {
	packageFile := c.Args().Get(0)
	workDir := c.Args().Get(1)
	// 如果路径存在目录或者文件
	if _, err := os.Stat(workDir); err == nil {
		return errors.Errorf("work directory (%s) path exisis", workDir)
	}
	// 创建目录
	if err := os.MkdirAll(workDir, 0775); err != nil {
		return err
	}
	if c.Bool("no-validate") {
		log.Println("[warn] package validation had been disabled!")
	}
	// 检查是否为题目包
	isDeerPack, err := utils2.IsProblemPackage(packageFile)
	if err != nil {
		return err
	}
	isZip, err := utils2.IsZipFile(packageFile)
	if err != nil {
		return err
	}
	// 解包
	if isDeerPack {
		if _, _, err := problems2.ReadProblemInfo(packageFile, true, !c.Bool("no-validate"), workDir); err != nil {
			return err
		}
	} else if isZip {
		if _, _, err := problems2.ReadProblemInfoZip(packageFile, true, !c.Bool("no-validate"), workDir); err != nil {
			return err
		}
		// clean meta file
		_ = os.Remove(path.Join(workDir, ".sign"))
		_ = os.Remove(path.Join(workDir, ".gpg"))
	} else {
		return errors.Errorf("not a deer-executor problem package file")
	}
	fmt.Println("Done.")
	return nil
}

// ReadProblemInfo 访问题目包信息
func ReadProblemInfo(c *cli.Context) error {
	packageFile := c.Args().Get(0)
	isDeerPack, err := utils2.IsProblemPackage(packageFile)
	if err != nil {
		return err
	}
	isZip, err := utils2.IsZipFile(packageFile)
	if err != nil {
		return err
	}
	// 如果是题目包文件，进行解包
	if isDeerPack {
		if c.Bool("gpg") {
			g, err := problems2.ReadProblemGPGInfo(packageFile)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
			fmt.Println(g)
		} else {
			s, _, err := problems2.ReadProblemInfo(packageFile, false, false, "")
			if err != nil {
				return err
			}
			fmt.Println(utils2.ObjectToJSONStringFormatted(s))
		}
	} else if isZip {
		if c.Bool("gpg") {
			g, err := problems2.ReadProblemGPGInfoZip(packageFile)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
			fmt.Println(g)
		} else {
			s, _, err := problems2.ReadProblemInfoZip(packageFile, false, false, "")
			if err != nil {
				return err
			}
			fmt.Println(utils2.ObjectToJSONStringFormatted(s))
		}
	} else {
		return errors.Errorf("not a deer-executor problem package file")
	}

	return nil
}
