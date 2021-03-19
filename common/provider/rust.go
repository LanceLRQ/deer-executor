package provider

import "fmt"

type RustCompileProvider struct {
	CodeCompileProvider
}

func NewRustCompileProvider() *RustCompileProvider {
	return &RustCompileProvider{
		CodeCompileProvider{
			isReady:  false,
			realTime: false,
			Name:     "rust",
		},
	}
}

func (prov *RustCompileProvider) Init(code string, workDir string) error {
	prov.codeContent = code
	prov.workDir = workDir

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".rs", "")
	return err
}

func (prov *RustCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.Rust, prov.codeFilePath, prov.programFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *RustCompileProvider) GetRunArgs() (args []string) {
	args = []string{prov.programFilePath}
	return
}

func (prov *RustCompileProvider) IsCompileError(remsg string) bool {
	return false
}
