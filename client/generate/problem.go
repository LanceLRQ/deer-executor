package generate

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/v3/executor"
	"github.com/LanceLRQ/deer-executor/v3/executor/persistence"
	"github.com/LanceLRQ/deer-executor/v3/executor/structs"
	utils "github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"path/filepath"
)

func makeProblmConfig() (*structs.JudgeConfiguration, error) {
	session, err := executor.NewSession("")
	if err != nil {
		return nil, err
	}
	config := session.JudgeConfig
	config.TestCases = []structs.TestCase{{}}
	config.Limitation = make(map[string]structs.JudgeResourceLimit)
	config.Limitation["gcc"] = structs.JudgeResourceLimit{
		TimeLimit:     config.TimeLimit,
		MemoryLimit:   config.MemoryLimit,
		RealTimeLimit: config.RealTimeLimit,
		FileSizeLimit: config.FileSizeLimit,
	}
	config.AnswerCases = []structs.AnswerCase{{}}
	config.SpecialJudge.CheckerCases = []structs.SpecialJudgeCheckerCase{{}}
	config.Problem.Sample = []structs.ProblemIOSample{{}}
	config.TestLib.ValidatorCases = []structs.TestlibValidatorCase{{}}
	config.TestLib.Generators = []structs.TestlibGenerator{{}}
	return &config, nil
}

// MakeProblemConfigFile Generate a problem configuration file.
func MakeProblemConfigFile(c *cli.Context) error {
	config, err := makeProblmConfig()
	if err != nil {
		return err
	}
	output := c.String("output")
	if output != "" {
		s, err := os.Stat(output)
		if s != nil || os.IsExist(err) {
			return errors.Errorf("output file exists")
		}
		fp, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return errors.Errorf("open output file error: %s\n", err.Error())
		}
		defer fp.Close()
		_, err = fp.WriteString(utils.ObjectToJSONStringFormatted(config))
		if err != nil {
			return err
		}
	} else {
		fmt.Println(utils.ObjectToJSONStringFormatted(config))
	}
	return nil
}

// InitProblemProjectDir create a problem project directory
func InitProblemProjectDir(c *cli.Context) error {
	workDir := c.Args().Get(0)
	// If the path exists and refers to a directory or a file.
	if _, err := os.Stat(workDir); err == nil {
		return errors.Errorf("work directory (%s) path exisis", workDir)
	}
	// make a new work directory
	if err := os.MkdirAll(workDir, 0775); err != nil {
		return err
	}
	example := c.String("sample")
	if example != "" {
		packageFile, err := filepath.Abs(example)
		if err != nil {
			return nil
		}
		// Check if the file belongs to deer-package
		yes, packageType, err := utils.IsDeerPackage(packageFile)
		if err != nil {
			return err
		}
		if !yes || packageType != persistence.PackageTypeProblem {
			return errors.Errorf("not a problem package")
		}
		pack, err := persistence.ParseProblemPackageFile(packageFile, true)
		if err != nil {
			return err
		}
		// unpack
		err = pack.UnpackProblemProject(workDir)
		if err != nil {
			return err
		}
	} else {
		// Create some commonly folders.
		dirs := []string{"answers", "cases", "bin", "codes", "generators"}
		for _, dirname := range dirs {
			err := os.MkdirAll(path.Join(workDir, dirname), 0775)
			if err != nil {
				return err
			}
		}
		// create a default conf
		config, err := makeProblmConfig()
		if err != nil {
			return err
		}
		// write to file
		if err = os.WriteFile(path.Join(workDir, "problem.json"), []byte(utils.ObjectToJSONStringFormatted(config)), 0664); err != nil {
			return err
		}
	}
	return nil
}
