package compile

import (
	"fmt"
)

type RubyCompileProvider struct {
	CodeCompileProvider
}


func (prov *RubyCompileProvider) Init(code string, workDir string) error {
	prov.isReady = false
	prov.realTime = true
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".rb", "")
	return err
}

func (prov *RubyCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(COMPILE_COMMAND_RUBY, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *RubyCompileProvider) GetRunArgs() (args []string) {
	args = []string{ "ruby", prov.codeFilePath }
	return
}

func (prov *RubyCompileProvider) IsCompileError(remsg string) bool {
	return false
}

