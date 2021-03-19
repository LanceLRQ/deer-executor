package provider

// Python Compiler Provider

import (
	"fmt"
	"strings"
)

// Py2CompileProvider python2语言编译提供程序
type Py2CompileProvider struct {
	CodeCompileProvider
}

// Py3CompileProvider python3语言编译提供程序
type Py3CompileProvider struct {
	CodeCompileProvider
}

// NewPy2CompileProvider 创建一个python2语言编译提供程序
func NewPy2CompileProvider() *Py2CompileProvider {
	return &Py2CompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: true,
			Name:     "python2",
		},
	}
}

// Init 初始化
func (prov *Py2CompileProvider) Init(code string, workDir string) error {
	prov.realTime = true
	prov.codeContent = code
	prov.workDir = workDir
	prov.Name = "python2"

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".py", "")
	if err != nil {
		return nil
	}
	prov.isReady = true
	return nil
}

// Compile 编译程序
func (prov *Py2CompileProvider) Compile() (result bool, errmsg string) {
	return true, ""
}

// GetRunArgs 获取运行参数
func (prov *Py2CompileProvider) GetRunArgs() (args []string) {
	argsRaw := fmt.Sprintf(CompileCommands.Python2, prov.codeFilePath)
	args = strings.Split(argsRaw, " ")
	return
}

// IsCompileError 是否编译错误
func (prov *Py2CompileProvider) IsCompileError(remsg string) bool {
	return strings.Contains(remsg, "SyntaxError") ||
		strings.Contains(remsg, "IndentationError") ||
		strings.Contains(remsg, "ImportError")
}

// NewPy3CompileProvider 创建一个python3语言编译提供程序
func NewPy3CompileProvider() *Py3CompileProvider {
	return &Py3CompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: true,
			Name:     "python3",
		},
	}
}

// Init 初始化
func (prov *Py3CompileProvider) Init(code string, workDir string) error {
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".py", "")
	if err != nil {
		return nil
	}
	prov.isReady = true
	return nil
}

// Compile 编译程序
func (prov *Py3CompileProvider) Compile() (result bool, errmsg string) {
	return true, ""
}

// GetRunArgs 获取运行参数
func (prov *Py3CompileProvider) GetRunArgs() (args []string) {
	argsRaw := fmt.Sprintf(CompileCommands.Python3, prov.codeFilePath)
	args = strings.Split(argsRaw, " ")
	return
}

// IsCompileError 是否编译错误
func (prov *Py3CompileProvider) IsCompileError(remsg string) bool {
	return strings.Contains(remsg, "SyntaxError") ||
		strings.Contains(remsg, "IndentationError") ||
		strings.Contains(remsg, "ImportError")
}
