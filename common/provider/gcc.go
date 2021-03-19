package provider

// GCC Compiler Provider

import "fmt"

// GnucCompileProvider c语言编译提供程序
type GnucCompileProvider struct {
	CodeCompileProvider
}

// GnucppCompileProvider c++语言编译提供程序
type GnucppCompileProvider struct {
	CodeCompileProvider
}

// NewGnucCompileProvider 创建一个c语言编译提供程序
func NewGnucCompileProvider() *GnucCompileProvider {
	return &GnucCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: false,
			Name:     "gcc",
		},
	}
}

// Init 初始化
func (prov *GnucCompileProvider) Init(code string, workDir string) error {
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".c", "")
	return err
}

// Compile 编译程序
func (prov *GnucCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.GNUC, prov.codeFilePath, prov.programFilePath))
	if result {
		prov.isReady = true
	}
	return
}

// GetRunArgs 获取运行参数
func (prov *GnucCompileProvider) GetRunArgs() (args []string) {
	args = []string{prov.programFilePath}
	return
}

// IsCompileError 是否编译错误
func (prov *GnucCompileProvider) IsCompileError(remsg string) bool {
	return false
}

// ManualCompile 执行手动编译
func (prov *GnucCompileProvider) ManualCompile(source string, target string, libraryDir []string) (bool, string) {
	cmd := fmt.Sprintf(CompileCommands.GNUC, source, target)
	if libraryDir != nil {
		for _, v := range libraryDir {
			cmd += fmt.Sprintf(" -I %s", v)
		}
	}
	result, err := prov.shell(cmd)
	return result, err
}

// NewGnucppCompileProvider 创建一个c++语言编译提供程序
func NewGnucppCompileProvider() *GnucppCompileProvider {
	return &GnucppCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: false,
			Name:     "g++",
		},
	}
}

// Init 初始化
func (prov *GnucppCompileProvider) Init(code string, workDir string) error {
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".cpp", "")
	return err
}

// Compile 编译程序
func (prov *GnucppCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.GNUCPP, prov.codeFilePath, prov.programFilePath))
	if result {
		prov.isReady = true
	}
	return
}

// GetRunArgs 获取运行参数
func (prov *GnucppCompileProvider) GetRunArgs() (args []string) {
	args = []string{prov.programFilePath}
	return
}

// IsCompileError 是否编译错误
func (prov *GnucppCompileProvider) IsCompileError(remsg string) bool {
	return false
}

// ManualCompile 执行手动编译
func (prov *GnucppCompileProvider) ManualCompile(source string, target string, libraryDir []string) (bool, string) {
	cmd := fmt.Sprintf(CompileCommands.GNUCPP, source, target)
	if libraryDir != nil {
		for _, v := range libraryDir {
			cmd += fmt.Sprintf(" -I %s", v)
		}
	}
	result, err := prov.shell(cmd)
	return result, err
}
