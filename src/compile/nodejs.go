package compile

import (
	"fmt"
	"strings"
)

type NodeJSCompileProvider struct {
	CodeCompileProvider
}


func (prov *NodeJSCompileProvider) Init(code string, workDir string) error {
	prov.isReady = false
	prov.realTime = true
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".js", "")
	return err
}

func (prov *NodeJSCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(COMPILE_COMMAND_NODEJS, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *NodeJSCompileProvider) GetRunArgs() (args []string) {
	args = []string{ "node", prov.codeFilePath }
	return
}

func (prov *NodeJSCompileProvider) IsCompileError(remsg string) bool {
	return strings.Contains(remsg, "SyntaxError") ||
		strings.Contains(remsg, "Error: Cannot find module")
}

