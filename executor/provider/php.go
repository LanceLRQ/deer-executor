package provider

// PHP Compiler Provider

import (
	"fmt"
)

// PHPCompileProvider php语言编译提供程序
type PHPCompileProvider struct {
	CodeCompileProvider
}

// NewPHPCompileProvider 创建一个php语言编译提供程序
func NewPHPCompileProvider() *PHPCompileProvider {
	return &PHPCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: true,
			Name:     "php",
		},
	}
}

// Init 初始化
func (prov *PHPCompileProvider) Init(code string, workDir string) error {
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".php", "")
	return err
}

// Compile 编译程序
func (prov *PHPCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.PHP, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

// GetRunArgs 获取运行参数
func (prov *PHPCompileProvider) GetRunArgs() (args []string) {
	args = []string{"/usr/bin/php", "-f", prov.codeFilePath}
	return
}

// IsCompileError 是否编译错误
func (prov *PHPCompileProvider) IsCompileError(remsg string) bool {
	return false
}
