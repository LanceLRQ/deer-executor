package compile

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

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".c", "")
	return err
}

func (prov *GnucCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(COMPILE_COMMAND_GNUC, prov.codeFilePath, prov.programFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *GnucCompileProvider) GetRunArgs() (args []string) {
	args = []string{ prov.programFilePath }
	return
}

func (prov *GnucppCompileProvider) Init(code string, workDir string) error {
	prov.isReady = false
	prov.realTime = false
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
	result, errmsg = prov.shell(fmt.Sprintf(COMPILE_COMMAND_GNUCPP, prov.codeFilePath, prov.programFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *GnucppCompileProvider) GetRunArgs() (args []string) {
	args = []string{ prov.programFilePath }
	return
}