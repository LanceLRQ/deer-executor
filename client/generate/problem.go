package generate

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/v3/executor"
	"github.com/LanceLRQ/deer-executor/v3/executor/persistence/problems"
	"github.com/LanceLRQ/deer-executor/v3/executor/structs"
	utils2 "github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path"
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

// MakeProblemConfigFile 生成评测配置文件
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
		_, err = fp.WriteString(utils2.ObjectToJSONStringFormatted(config))
		if err != nil {
			return err
		}
	} else {
		fmt.Println(utils2.ObjectToJSONStringFormatted(config))
	}
	return nil
}

// InitProblemWorkDir 创建一个题目工作目录
func InitProblemWorkDir(c *cli.Context) error {
	workDir := c.Args().Get(0)
	// 如果路径存在目录或者文件
	if _, err := os.Stat(workDir); err == nil {
		return errors.Errorf("work directory (%s) path exisis", workDir)
	}
	// 创建目录
	if err := os.MkdirAll(workDir, 0775); err != nil {
		return err
	}
	example := c.String("name")
	if example != "" {
		packageFile := path.Join("./lib/example", example)
		// 检查题目包是否存在
		yes, err := utils2.IsProblemPackage(packageFile)
		if err != nil {
			return err
		}
		if !yes {
			return errors.Errorf("not a problem package")
		}
		// 如果指定了对应的模板
		if _, _, err := problems.ReadProblemInfo(packageFile, true, true, workDir); err != nil {
			return err
		}
	} else {
		// 创建文件夹
		dirs := []string{"answers", "cases", "bin", "codes", "generators"}
		for _, dirname := range dirs {
			err := os.MkdirAll(path.Join(workDir, dirname), 0775)
			if err != nil {
				return err
			}
		}
		/// 创建配置
		config, err := makeProblmConfig()
		if err != nil {
			return err
		}
		// 写入到文件
		if err = ioutil.WriteFile(path.Join(workDir, "problem.json"), []byte(utils2.ObjectToJSONStringFormatted(config)), 0664); err != nil {
			return err
		}
	}
	return nil
}
