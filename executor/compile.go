package executor

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/provider"
	"path"
	"strings"
)

// 匹配编程语言
func matchCodeLanguage(keyword string, fileName string) (provider.CodeCompileProviderInterface, error) {
_match:
	switch keyword {
	case "c", "gcc", "gnu-c":
		return &provider.GnucCompileProvider{}, nil
	case "cpp", "gcc-cpp", "gpp", "g++":
		return &provider.GnucppCompileProvider{}, nil
	case "java":
		return &provider.JavaCompileProvider{}, nil
	case "py2", "python2":
		return &provider.Py2CompileProvider{}, nil
	case "py", "py3", "python3":
		return &provider.Py3CompileProvider{}, nil
	case "php":
		return &provider.PHPCompileProvider{}, nil
	case "node", "nodejs":
		return &provider.NodeJSCompileProvider{}, nil
	case "rb", "ruby":
		return &provider.RubyCompileProvider{}, nil
	case "auto":
		keyword = strings.Replace(path.Ext(fileName), ".", "", -1)
		goto _match
	}
	return nil, fmt.Errorf("unsupported language")
}

// 编译文件
// 如果不设置codeStr，默认会读取配置文件里的code_file字段并打开对应文件
func (session *JudgeSession) getCompiler(codeStr string) (provider.CodeCompileProviderInterface, error) {
	if codeStr == "" {
		codeFileBytes, err := ReadFile(session.CodeFile)
		if err != nil {
			return nil, err
		}
		codeStr = string(codeFileBytes)
	}

	compiler, err := matchCodeLanguage(session.CodeLangName, session.CodeFile)
	if err != nil { return nil, err }
	err = compiler.Init(codeStr, "/tmp")
	if err != nil {
		return nil, err
	}
	return compiler, err
}

// 编译目标程序
func (session *JudgeSession)compileTargetProgram(judgeResult *JudgeResult) error {
	// 获取对应的编译器提供程序
	compiler, err := session.getCompiler("")
	if err != nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = err.Error()
		return err
	}
	// 编译程序
	success, ceinfo := compiler.Compile()
	if !success {
		judgeResult.JudgeResult = JudgeFlagCE
		judgeResult.CeInfo = ceinfo
		return fmt.Errorf("compile error:\n%s", ceinfo)
	}
	// 清理工作目录
	defer compiler.Clean()
	// 获取执行指令
	session.Commands = compiler.GetRunArgs()
	return nil
}