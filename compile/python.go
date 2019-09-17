/* Python Compiler Provider
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPL.
 */
package deer_compiler

import "strings"

type Py2CompileProvider struct {
	CodeCompileProvider
}

type Py3CompileProvider struct {
	CodeCompileProvider
}

func (prov *Py2CompileProvider) Init(code string, workDir string) error {
	prov.realTime = true
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".py", "")
	if err != nil {
		return nil
	}
	prov.isReady = true
	return nil
}

func (prov *Py2CompileProvider) Compile() (result bool, errmsg string) {
	return true, ""
}

func (prov *Py2CompileProvider) GetRunArgs() (args []string) {
	args = []string{ "/usr/bin/python", "-u", prov.codeFilePath }
	return
}

func (prov *Py2CompileProvider) IsCompileError(remsg string) bool {
	return strings.Contains(remsg, "SyntaxError") ||
		strings.Contains(remsg, "IndentationError") ||
		strings.Contains(remsg, "ImportError")
}


func (prov *Py3CompileProvider) Init(code string, workDir string) error {
	prov.isReady = false
	prov.realTime = true
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".py", "")
	if err != nil {
		return nil
	}
	prov.isReady = true
	return nil
}

func (prov *Py3CompileProvider) Compile() (result bool, errmsg string) {
	return true, ""
}

func (prov *Py3CompileProvider) GetRunArgs() (args []string) {
	args = []string{ "/usr/bin/python3", "-u", prov.codeFilePath }
	return
}

func (prov *Py3CompileProvider) IsCompileError(remsg string) bool {
	return strings.Contains(remsg, "SyntaxError") ||
		strings.Contains(remsg, "IndentationError") ||
		strings.Contains(remsg, "ImportError")
}