package provider

// Ruby Compiler Provider

import (
	"fmt"
)

// RubyCompileProvider ruby语言编译提供程序
type RubyCompileProvider struct {
	CodeCompileProvider
}

// NewRubyCompileProvider 创建一个ruby语言编译提供程序
func NewRubyCompileProvider() *RubyCompileProvider {
	return &RubyCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: true,
			Name:     "ruby",
		},
	}
}

// Init 初始化
func (prov *RubyCompileProvider) Init(code string, workDir string) error {
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".rb", "")
	return err
}

// Compile 编译程序
func (prov *RubyCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.Ruby, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

// GetRunArgs 获取运行参数
func (prov *RubyCompileProvider) GetRunArgs() (args []string) {
	args = []string{"/usr/bin/ruby", prov.codeFilePath}
	return
}

// IsCompileError 是否编译错误
func (prov *RubyCompileProvider) IsCompileError(remsg string) bool {
	return false
}
