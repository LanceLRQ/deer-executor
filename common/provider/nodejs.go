package provider

// NodeJS Compiler Provider

import (
	"fmt"
	"strings"
)

// NodeJSCompileProvider nodejs语言编译提供程序
type NodeJSCompileProvider struct {
	CodeCompileProvider
}

// NewNodeJSCompileProvider 创建一个nodejs语言编译提供程序
func NewNodeJSCompileProvider() *NodeJSCompileProvider {
	return &NodeJSCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: true,
			Name:     "nodejs",
		},
	}
}

// Init 初始化
func (prov *NodeJSCompileProvider) Init(code string, workDir string) error {
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".js", "")
	return err
}

// Compile 编译程序
func (prov *NodeJSCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.NodeJS, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

// GetRunArgs 获取运行参数
func (prov *NodeJSCompileProvider) GetRunArgs() (args []string) {
	args = []string{"/usr/bin/node", prov.codeFilePath}
	return
}

// IsCompileError 是否编译错误
func (prov *NodeJSCompileProvider) IsCompileError(remsg string) bool {
	return strings.Contains(remsg, "SyntaxError") ||
		strings.Contains(remsg, "Error: Cannot find module")
}
