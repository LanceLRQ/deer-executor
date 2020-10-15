/* Java Compiler Provider
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package provider

import (
	"fmt"
	"path"
	"regexp"
)

type JavaCompileProvider struct {
	CodeCompileProvider
	javaClassName string
}

func getJavaClassName (code string) (className string, err error) {
	reg := regexp.MustCompile(`public class ([A-Za-z0-9_$]+)`)
	matched := reg.FindSubmatch([]byte(code))
	if matched != nil {
		className = string(matched[1])
		err = nil
	} else {
		className = "Main"			// default java public classname (might cause compile error)
 		// err = fmt.Errorf("no java public class name matched")
	}
	return
}

func (prov *JavaCompileProvider) Init(code string, workDir string) error {
	prov.isReady = false
	prov.realTime = false
	prov.codeContent = code
	prov.workDir = workDir
	prov.Name = "Java"

	javaClassName, err := getJavaClassName(prov.codeContent)
	if err != nil {
		return err
	}
	prov.javaClassName = javaClassName

	err = prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".java", ".class")
	return err
}

func (prov *JavaCompileProvider) initFiles(codeExt string, programExt string) error {
	prov.codeFileName = fmt.Sprintf("%s%s", prov.javaClassName, codeExt)
	prov.programFileName = fmt.Sprintf("%s%s", prov.javaClassName, programExt)
	prov.codeFilePath = path.Join(prov.workDir, prov.codeFileName)
	prov.programFilePath = path.Join(prov.workDir, prov.programFileName)
	err := prov.saveCode()
	return err
}


func (prov *JavaCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommandJava, prov.codeFilePath, path.Dir(prov.programFilePath)))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *JavaCompileProvider) GetRunArgs() (args []string) {
	args = []string{
		"/usr/bin/java", "-client", "-Dfile.encoding=utf-8",
		"-classpath", path.Dir(prov.programFilePath), prov.javaClassName,
	}
	return
}

func (prov *JavaCompileProvider) IsCompileError(remsg string) bool {
	return false
}