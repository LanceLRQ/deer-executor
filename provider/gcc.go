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

func (prov *GnucCompileProvider) Init(code string, workDir string) error {
	prov.isReady = false
	prov.realTime = false
	prov.codeContent = code
	prov.workDir = workDir
	prov.Name = "GCC"

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

func (prov *GnucCompileProvider) GetRunArgs() (args []string) {
	args = []string{ prov.programFilePath }
	return
}

func (prov *GnucCompileProvider) IsCompileError(remsg string) bool {
	return false
}

func (prov *GnucppCompileProvider) Init(code string, workDir string) error {
	prov.isReady = false
	prov.realTime = false
	prov.codeContent = code
	prov.workDir = workDir
	prov.Name = "GCC-CPP"

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
	args = []string{ prov.programFilePath }
	return
}

func (prov *GnucppCompileProvider) IsCompileError(remsg string) bool {
	return false
}