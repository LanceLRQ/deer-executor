package compile

import "fmt"

type GolangCompileProvider struct {
	CodeCompileProvider
}


func (prov *GolangCompileProvider) Init(code string, workDir string) error {
	prov.IsReady = false
	prov.RealTime = false
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
	result, errmsg = prov.shell(fmt.Sprintf(COMPILE_COMMAND_GO, prov.programFilePath, prov.codeFilePath))
	if result {
		prov.IsReady = true
	}
	return
}

func (prov *GolangCompileProvider) GetRunArgs() (args []string) {
	args = []string{ prov.programFilePath }
	return
}
