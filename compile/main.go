/* Compiler Provider Base
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package deer_compiler

import (
	"bytes"
	"fmt"
	"github.com/satori/go.uuid"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	COMPILE_COMMAND_GNUC = "/usr/bin/gcc %s -o %s -ansi -fno-asm -Wall -std=c99 -lm"
	COMPILE_COMMAND_GNUCPP = "/usr/bin/g++ %s -o %s -ansi -fno-asm -Wall -lm"
	COMPILE_COMMAND_JAVA = "/usr/bin/javac -encoding utf-8 %s -d %s"
	COMPILE_COMMAND_GO = "/usr/bin/go build -o %s %s"
	COMPILE_COMMAND_NODEJS = "/usr/bin/node -c %s"
	COMPILE_COMMAND_PHP = "/usr/bin/php -l -f %s"
	COMPILE_COMMAND_RUBY = "/usr/bin/ruby -c %s"
)

type CodeCompileProviderInterface interface {
	// 初始化
	Init(code string, workDir string) error
	// 初始化文件信息
	initFiles(codeExt string, programExt string) error
	// 编译程序
	Compile() (result bool, errmsg string)
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
}

type CodeCompileProvider struct {
	CodeCompileProviderInterface
	codeContent string		// 代码
	realTime bool			// 是否为实时编译的语言
	isReady bool			// 是否已经编译完毕
	codeFileName, codeFilePath string			// 目标程序源文件
	programFileName, programFilePath string		// 目标程序文件
	workDir string			// 工作目录
}

func (prov *CodeCompileProvider) initFiles(codeExt string, programExt string) error {
	prov.codeFileName = fmt.Sprintf("%s%s", uuid.NewV4().String(), codeExt)
	prov.programFileName = fmt.Sprintf("%s%s", uuid.NewV4().String(), programExt)
	prov.codeFilePath = path.Join(prov.workDir, prov.codeFileName)
	prov.programFilePath = path.Join(prov.workDir, prov.programFileName)

	err := prov.saveCode()
	return err
}

func (prov *CodeCompileProvider) shell(commands string) (success bool, errout string) {
	cmdArgs := strings.Split(commands, " ")
	if len(cmdArgs) <= 1 {
		return false, "Not enough arguments for compiler"
	}
	proc := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	var stderr bytes.Buffer
	proc.Stderr = &stderr
	err := proc.Run()
	if err != nil {
		return false, stderr.String()
	}
	return true, ""
}

func (prov *CodeCompileProvider) saveCode() error {
	file, err := os.OpenFile(prov.codeFilePath, os.O_RDWR | os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.WriteString(prov.codeContent)
	return err
}

func (prov *CodeCompileProvider) checkWorkDir() error {
	_, err := os.Stat(prov.workDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("work dir not exists")
		}
		return err
	}
	return nil
}

func (prov *CodeCompileProvider) IsRealTime() bool {
	return prov.realTime
}

func (prov *CodeCompileProvider) IsReady() bool {
	return prov.isReady
}