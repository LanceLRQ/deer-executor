package packmgr

import (
	"context"
	"fmt"
	"github.com/LanceLRQ/deer-executor/v2/common/structs"
	"github.com/LanceLRQ/deer-executor/v2/common/utils"
	"github.com/LanceLRQ/deer-executor/v2/executor"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// 运行测试数据生成
func runTestCaseGen(session *executor.JudgeSession, tCase *structs.TestCase, withAnswer bool) error {
	// 如果是generator脚本
	if tCase.UseGenerator {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		inbytes, err := utils.CallGenerator(ctx, tCase, session.ConfigDir)
		if err != nil {
			return err
		}
		// 写入到文件
		err = ioutil.WriteFile(path.Join(session.ConfigDir, tCase.Input), inbytes, 0664)
		if err != nil {
			return err
		}
		log.Printf("[generator] generate input done!")
	}
	if withAnswer {
		fin, err := os.Open(path.Join(session.ConfigDir, tCase.Input))
		if err != nil {
			return err
		}
		fout, err := os.Create(path.Join(session.ConfigDir, tCase.Output))
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		rel, err := utils.RunUnixShell(&structs.ShellOptions{
			Context: ctx,
			Name:    session.Commands[0],
			Args:    session.Commands[1:],
			StdWriter: &structs.ShellWriters{
				Input:  fin,
				Output: fout,
				Error:  nil,
			},
		})
		if err != nil {
			return err
		}
		if !rel.Success {
			log.Printf("[generator] run answer code error: %s", rel.Stderr)
			return errors.Errorf("[generator] run answer code error: %s", rel.Stderr)
		}
		log.Printf("[generator] generate answer done!")
	}
	return nil
}

// 运行test cases的数据生成
// caseIndex < 0 表示校验全部
func runTestCaseGenerator(session *executor.JudgeSession, caseIndex int, withAnswer bool) error {
	// 执行遍历
	if caseIndex < 0 {
		for key := range session.JudgeConfig.TestCases {
			log.Printf("[generator] run case #%d", key)
			err := runTestCaseGen(session, &session.JudgeConfig.TestCases[key], withAnswer)
			if err != nil {
				return err
			}
		}
	} else {
		log.Printf("[generator] run case #%d", caseIndex)
		err := runTestCaseGen(session, &session.JudgeConfig.TestCases[caseIndex], withAnswer)
		if err != nil {
			return err
		}
	}
	return nil
}

func initWork(session *executor.JudgeSession, answerCaseIndex uint) error {

	// 强制设定工作目录
	session.SessionID = uuid.NewV4().String()
	session.SessionRoot = "/tmp"
	// 初始化session dir
	sessionDir, err := utils.GetSessionDir(session.SessionRoot, session.SessionID)
	if err != nil {
		return err
	}
	session.SessionDir = sessionDir

	acase := session.JudgeConfig.AnswerCases[answerCaseIndex]
	// 如果有代码文件
	if acase.FileName != "" {
		session.CodeFile = path.Join(session.ConfigDir, acase.FileName)
	}
	session.CodeLangName = acase.Language
	// 获取对应的编译器提供程序
	compiler, err := session.GetCompiler(acase.Content)
	if err != nil {
		return err
	}
	// 编译程序
	success, ceinfo := compiler.Compile()
	if !success {
		return errors.Errorf("[generator] compile error:\n%s", ceinfo)
	}
	// 获取执行指令
	session.Commands = compiler.GetRunArgs()
	session.Compiler = compiler
	return nil
}

// RunTestCaseGenerator 运行测试数据生成器 (APP入口)
func RunTestCaseGenerator(c *cli.Context) error {
	configFile := c.Args().Get(0)
	_, err := os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		return errors.Errorf("[generator] problem config file (%s) not found", configFile)
	}
	session, err := executor.NewSession(configFile)
	if err != nil {
		return err
	}

	if !c.Bool("silence") {
		fmt.Print("[generator] Operation will cover all files, continue? [y/N] ")
		ans := ""
		_, err := fmt.Scanf("%s", &ans)
		if err != nil {
			return nil // don't crash at EOF
		}
		if len(ans) > 0 && strings.ToLower(ans[:1]) != "y" {
			return nil
		}
	}

	withAnswer := c.Bool("with-answer")
	answerCaseIndex := c.Uint("answer")
	testCaseIndex := c.Int("case")

	// 编译答案代码
	if withAnswer {
		err = initWork(session, answerCaseIndex)
		if err != nil {
			return err
		}
		defer session.Clean()
	}

	return runTestCaseGenerator(session, testCaseIndex, withAnswer)
}
