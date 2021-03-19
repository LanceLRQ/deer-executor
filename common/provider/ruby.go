/* Ruby Compiler Provider
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package provider

import (
	"fmt"
)

type RubyCompileProvider struct {
	CodeCompileProvider
}

func NewRubyCompileProvider() *RubyCompileProvider {
	return &RubyCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: true,
			Name:     "ruby",
		},
	}
}

func (prov *RubyCompileProvider) Init(code string, workDir string) error {
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
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.Ruby, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *RubyCompileProvider) GetRunArgs() (args []string) {
	args = []string{"/usr/bin/ruby", prov.codeFilePath}
	return
}

func (prov *RubyCompileProvider) IsCompileError(remsg string) bool {
	return false
}
