package provider

// Golang Compiler Provider

import "fmt"

// GolangCompileProvider go语言编译提供程序
type GolangCompileProvider struct {
	CodeCompileProvider
}

// NewGolangCompileProvider 创建一个go语言编译提供程序
func NewGolangCompileProvider() *GolangCompileProvider {
	return &GolangCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: false,
			Name:     "golang",
		},
	}
}

// Init 初始化
func (prov *GolangCompileProvider) Init(code string, workDir string) error {
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".go", "")
	return err
}

// Compile 编译程序
func (prov *GolangCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.Go, prov.programFilePath, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

// GetRunArgs 获取运行参数
func (prov *GolangCompileProvider) GetRunArgs() (args []string) {
	args = []string{prov.programFilePath}
	return
}

// IsCompileError 是否编译错误
func (prov *GolangCompileProvider) IsCompileError(remsg string) bool {
	return false
}

// ManualCompile 执行手动编译
func (prov *GolangCompileProvider) ManualCompile(source string, target string) (bool, string) {
	cmd := fmt.Sprintf(CompileCommands.Go, source, target)
	result, err := prov.shell(cmd)
	return result, err
}
