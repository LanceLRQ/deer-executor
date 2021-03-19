package provider

import "fmt"

// RustCompileProvider rust语言编译提供程序
type RustCompileProvider struct {
	CodeCompileProvider
}

// NewRustCompileProvider 创建一个rust语言编译提供程序
func NewRustCompileProvider() *RustCompileProvider {
	return &RustCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: false,
			Name:     "rust",
		},
	}
}

// Init 初始化
func (prov *RustCompileProvider) Init(code string, workDir string) error {
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".rs", "")
	return err
}

// Compile 编译程序
func (prov *RustCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.Rust, prov.codeFilePath, prov.programFilePath))
	if result {
		prov.isReady = true
	}
	return
}

// GetRunArgs 获取运行参数
func (prov *RustCompileProvider) GetRunArgs() (args []string) {
	args = []string{prov.programFilePath}
	return
}

// IsCompileError 是否编译错误
func (prov *RustCompileProvider) IsCompileError(remsg string) bool {
	return false
}
