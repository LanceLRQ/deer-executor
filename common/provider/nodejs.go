/* NodeJS Compiler Provider
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package provider

import (
	"fmt"
	"strings"
)

type NodeJSCompileProvider struct {
	CodeCompileProvider
}

func NewNodeJSCompileProvider() *NodeJSCompileProvider {
	return &NodeJSCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: true,
			Name:     "nodejs",
		},
	}
}

func (prov *NodeJSCompileProvider) Init(code string, workDir string) error {
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
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.NodeJS, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *NodeJSCompileProvider) GetRunArgs() (args []string) {
	args = []string{"/usr/bin/node", prov.codeFilePath}
	return
}

func (prov *NodeJSCompileProvider) IsCompileError(remsg string) bool {
	return strings.Contains(remsg, "SyntaxError") ||
		strings.Contains(remsg, "Error: Cannot find module")
}
