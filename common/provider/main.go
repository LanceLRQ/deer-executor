package provider

// Compiler Provider Base

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LanceLRQ/deer-executor/v2/common/structs"
	"github.com/LanceLRQ/deer-executor/v2/common/utils"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

// CompileCommandsStruct 编译命令集
type CompileCommandsStruct struct {
	GNUC    string `json:"gcc"`
	GNUCPP  string `json:"g++"`
	Java    string `json:"java"`
	Go      string `json:"golang"`
	NodeJS  string `json:"nodejs"`
	PHP     string `json:"php"`
	Ruby    string `json:"ruby"`
	Python2 string `json:"python2"`
	Python3 string `json:"python3"`
	Rust    string `json:"rust"`
}

// CompileCommands 定义默认的编译命令集
var CompileCommands = CompileCommandsStruct{
	GNUC:    "gcc %s -o %s -ansi -fno-asm -Wall -std=c11 -lm",
	GNUCPP:  "g++ %s -o %s -ansi -fno-asm -Wall -lm -std=c++11",
	Java:    "javac -encoding utf-8 %s -d %s",
	Go:      "go build -o %s %s",
	NodeJS:  "node -c %s",
	PHP:     "php -l -f %s",
	Ruby:    "ruby -c %s",
	Python2: "python -u %s",
	Python3: "python3 -u %s",
	Rust:    "rustc %s -o %s",
}

// CodeCompileProviderInterface 代码编译提供程序接口定义
type CodeCompileProviderInterface interface {
	// 初始化
	Init(code string, workDir string) error
	// 初始化文件信息
	initFiles(codeExt string, programExt string) error
	// 执行编译
	Compile() (result bool, errmsg string)
	// 清理工作目录
	Clean()
	// 获取程序的运行命令参数组
	GetRunArgs() (args []string)
	// 判断STDERR的输出内容是否存在编译错误信息，通常用于脚本语言的判定，
	IsCompileError(remsg string) bool
	// 是否为实时编译的语言
	IsRealTime() bool
	// 是否已经编译完毕
	IsReady() bool
	// 调用Shell命令并获取运行结果
	shell(commands string) (success bool, errout string)
	// 保存代码到文件
	saveCode() error
	// 检查工作目录是否存在
	checkWorkDir() error
	// 获取提供程序的名称
	GetName() string
}

// CodeCompileProvider 代码编译提供程序公共结构定义
type CodeCompileProvider struct {
	CodeCompileProviderInterface
	Name                             string // 编译器提供程序名称
	codeContent                      string // 代码
	realTime                         bool   // 是否为实时编译的语言
	isReady                          bool   // 是否已经编译完毕
	codeFileName, codeFilePath       string // 目标程序源文件
	programFileName, programFilePath string // 目标程序文件
	workDir                          string // 工作目录
}

// PlaceCompilerCommands 替换编译命令集
func PlaceCompilerCommands(configFile string) error {
	if configFile != "" {
		_, err := os.Stat(configFile)
		// ignore
		if os.IsNotExist(err) {
			return nil
		}
		cbody, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(cbody, &CompileCommands)
		if err != nil {
			return err
		}
	}
	return nil
}

// 初始化文件
func (prov *CodeCompileProvider) initFiles(codeExt string, programExt string) error {
	prov.codeFileName = fmt.Sprintf("%s%s", uuid.NewV4().String(), codeExt)
	prov.programFileName = fmt.Sprintf("%s%s", uuid.NewV4().String(), programExt)
	prov.codeFilePath = path.Join(prov.workDir, prov.codeFileName)
	prov.programFilePath = path.Join(prov.workDir, prov.programFileName)

	err := prov.saveCode()
	return err
}

// GetName 获取提供程序的名称
func (prov *CodeCompileProvider) GetName() string {
	return prov.Name
}

// Clean 清理代码
func (prov *CodeCompileProvider) Clean() {
	_ = os.Remove(prov.codeFilePath)
	_ = os.Remove(prov.programFilePath)
}

// 执行shell
func (prov *CodeCompileProvider) shell(commands string) (success bool, errout string) {
	ctx, _ := context.WithTimeout(context.Background(), 7*time.Second)
	cmdArgs := strings.Split(commands, " ")
	if len(cmdArgs) <= 1 {
		return false, "not enough arguments for compiler"
	}
	ret, err := utils.RunUnixShell(&structs.ShellOptions{
		Context:   ctx,
		Name:      cmdArgs[0],
		Args:      cmdArgs[1:],
		StdWriter: nil,
		OnStart:   nil,
	})
	if err != nil {
		return false, err.Error()
	}
	if !ret.Success {
		return false, ret.Stderr
	}
	return true, ""
}

// 存储代码到文件
func (prov *CodeCompileProvider) saveCode() error {
	file, err := os.OpenFile(prov.codeFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(prov.codeContent)
	return err
}

// 检查工作目录
func (prov *CodeCompileProvider) checkWorkDir() error {
	_, err := os.Stat(prov.workDir)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Errorf("work dir not exists")
		}
		return err
	}
	return nil
}

// IsRealTime 返回是否是实时评测语言
func (prov *CodeCompileProvider) IsRealTime() bool {
	return prov.realTime
}

// IsReady 返回编译是否就绪
func (prov *CodeCompileProvider) IsReady() bool {
	return prov.isReady
}
