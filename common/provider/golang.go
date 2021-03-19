/* Golang Compiler Provider
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package provider

import "fmt"

type GolangCompileProvider struct {
	CodeCompileProvider
}

func NewGolangCompileProvider() *GolangCompileProvider {
	return &GolangCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: false,
			Name:     "golang",
		},
	}
}

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

func (prov *GolangCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.Go, prov.programFilePath, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *GolangCompileProvider) GetRunArgs() (args []string) {
	args = []string{prov.programFilePath}
	return
}

func (prov *GolangCompileProvider) IsCompileError(remsg string) bool {
	return false
}

// 手动编译
func (prov *GolangCompileProvider) ManualCompile(source string, target string) (bool, string) {
	cmd := fmt.Sprintf(CompileCommands.Go, source, target)
	result, err := prov.shell(cmd)
	return result, err
}
