package packmgr

import (
	"context"
	"github.com/LanceLRQ/deer-executor/v3/executor"
	structs2 "github.com/LanceLRQ/deer-executor/v3/executor/structs"
	utils2 "github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func runValidatorCase(vBin string, vCase *structs2.TestlibValidatorCase) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	rel, err := utils2.RunUnixShell(&structs2.ShellOptions{
		Context:   ctx,
		Name:      vBin,
		Args:      nil,
		StdWriter: nil,
		OnStart: func(writer io.Writer) error {
			_, err := writer.Write([]byte(vCase.Input))
			return err
		},
	})
	if err != nil {
		return err
	}
	if rel.Success {
		vCase.ValidatorVerdict = true
		vCase.ValidatorComment = ""
	} else {
		log.Printf("[validator] validator error: %s", rel.Stderr)
		vCase.ValidatorVerdict = false
		vCase.ValidatorComment = rel.Stderr
	}
	vCase.Verdict = vCase.ValidatorVerdict == vCase.ExpectedVerdict
	return nil
}

func runTestCase(configDir, vBin string, tCase *structs2.TestCase) error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	var inbytes []byte
	var err error
	// 判断是generator还是普通input
	if tCase.UseGenerator {
		inbytes, err = utils2.CallGenerator(ctx, tCase, configDir)
		if err != nil {
			return err
		}
	} else {
		inbytes, err = ioutil.ReadFile(path.Join(configDir, tCase.Input))
		if err != nil {
			return err
		}
	}

	rel, err := utils2.RunUnixShell(&structs2.ShellOptions{
		Context:   ctx,
		Name:      vBin,
		Args:      nil,
		StdWriter: nil,
		OnStart: func(writer io.Writer) error {
			_, err := writer.Write(inbytes)
			return err
		},
	})
	if err != nil {
		return err
	}
	if rel.Success {
		tCase.ValidatorVerdict = true
		tCase.ValidatorComment = ""
	} else {
		log.Printf("[validator] validator error: %s", rel.Stderr)
		tCase.ValidatorVerdict = false
		tCase.ValidatorComment = rel.Stderr
	}
	return nil
}

func isValidatorExists(config *structs2.JudgeConfiguration) error {
	validator, err := utils2.GetCompiledBinaryFileAbsPath("validator", config.TestLib.ValidatorName, config.ConfigDir)
	if err != nil {
		return err
	}
	_, err = os.Stat(validator)
	if os.IsNotExist(err) {
		return errors.Errorf("[validator] execuable validator file not exists")
	}
	return nil
}

// 运行Testlib的validator校验
func runTestlibValidators(config *structs2.JudgeConfiguration, moduleName string, caseIndex int) error {
	if err := isValidatorExists(config); err != nil {
		return err
	}
	if moduleName == "all" {
		caseIndex = -1
	}
	if moduleName == "all" || moduleName == "validator_cases" {
		if err := RunTestlibValidatorCases(config, caseIndex); err != nil {
			return err
		}
	}
	// for test_cases
	if moduleName == "all" || moduleName == "test_cases" {
		if err := RunTestCasesInputValidation(config, -1); err != nil {
			return err
		}
	}
	return nil
}

// RunTestlibValidatorCases 运行validator cases的校验
// caseIndex < 0 表示校验全部
func RunTestlibValidatorCases(config *structs2.JudgeConfiguration, caseIndex int) error {
	validator, err := utils2.GetCompiledBinaryFileAbsPath("validator", config.TestLib.ValidatorName, config.ConfigDir)
	if err != nil {
		return err
	}
	// 执行遍历
	if caseIndex < 0 {
		for key := range config.TestLib.ValidatorCases {
			log.Printf("[validator] run validator case #%d", key)
			err := runValidatorCase(validator, &config.TestLib.ValidatorCases[key])
			if err != nil {
				return err
			}
		}
	} else {
		log.Printf("[validator] run validator case #%d", caseIndex)
		err := runValidatorCase(validator, &config.TestLib.ValidatorCases[caseIndex])
		if err != nil {
			return err
		}
	}
	return nil
}

// RunTestCasesInputValidation 运行test cases的校验
// caseIndex < 0 表示校验全部
func RunTestCasesInputValidation(config *structs2.JudgeConfiguration, caseIndex int) error {
	validator, err := utils2.GetCompiledBinaryFileAbsPath("validator", config.TestLib.ValidatorName, config.ConfigDir)
	if err != nil {
		return err
	}
	// 执行遍历
	if caseIndex < 0 {
		for key := range config.TestCases {
			log.Printf("[validator] run test case #%d", key)
			err := runTestCase(config.ConfigDir, validator, &config.TestCases[key])
			if err != nil {
				return err
			}
		}
	} else {
		log.Printf("[validator] run test case #%d", caseIndex)
		err := runTestCase(config.ConfigDir, validator, &config.TestCases[caseIndex])
		if err != nil {
			return err
		}
	}
	return nil
}

// RunTestlibValidators 运行Testlib的validator校验 (APP入口)
func RunTestlibValidators(c *cli.Context) error {
	configFile := c.Args().Get(0)
	_, err := os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		return errors.Errorf("[validator] problem config file (%s) not found", configFile)
	}
	session, err := executor.NewSession(configFile)
	if err != nil {
		return err
	}
	mtype := c.String("type")
	mCaseIndex := c.Int("case")
	silence := c.Bool("silence")

	LIST := []string{"all", "validator_cases", "test_cases"}
	if !utils2.Contains(LIST, mtype) {
		return errors.Errorf("[validator] unsupport module type")
	}
	err = runTestlibValidators(&session.JudgeConfig, mtype, mCaseIndex)
	if err != nil {
		return err
	}
	return session.SaveConfiguration(!silence)
}
