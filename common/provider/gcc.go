/* GCC Compiler Provider
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package provider

import "fmt"

type GnucCompileProvider struct {
	CodeCompileProvider
}

type GnucppCompileProvider struct {
	CodeCompileProvider
}

func NewGnucCompileProvider() *GnucCompileProvider {
	return &GnucCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: false,
			Name:     "gcc",
		},
	}
}

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

func (prov *GnucCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.GNUC, prov.codeFilePath, prov.programFilePath))
	if result {
		prov.isReady = true
	}
	return
}

// 手动编译
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

func (prov *GnucCompileProvider) GetRunArgs() (args []string) {
	args = []string{prov.programFilePath}
	return
}

func (prov *GnucCompileProvider) IsCompileError(remsg string) bool {
	return false
}

func NewGnucppCompileProvider() *GnucppCompileProvider {
	return &GnucppCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: false,
			Name:     "g++",
		},
	}
}

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

func (prov *GnucppCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.GNUCPP, prov.codeFilePath, prov.programFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *GnucppCompileProvider) GetRunArgs() (args []string) {
	args = []string{prov.programFilePath}
	return
}

func (prov *GnucppCompileProvider) IsCompileError(remsg string) bool {
	return false
}

// 手动编译
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
