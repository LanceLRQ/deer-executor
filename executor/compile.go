package executor

import (
	"fmt"
	commonStructs "github.com/LanceLRQ/deer-common/structs"
	"github.com/LanceLRQ/deer-executor/provider"
	"io/ioutil"
	"os"
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
	case "go", "golang":
		return &provider.GolangCompileProvider{}, nil
	case "node", "nodejs":
		return &provider.NodeJSCompileProvider{}, nil
	case "rb", "ruby":
		return &provider.RubyCompileProvider{}, nil
	case "auto", "":
		keyword = strings.Replace(path.Ext(fileName), ".", "", -1)
		goto _match
	}
	return nil, fmt.Errorf("unsupported language")
}

// 编译文件
// 如果不设置codeStr，默认会读取配置文件里的code_file字段并打开对应文件
func (session *JudgeSession) getCompiler(codeStr string) (provider.CodeCompileProviderInterface, error) {
	if codeStr == "" {
		codeFileBytes, err := ioutil.ReadFile(session.CodeFile)
		if err != nil {
			return nil, err
		}
		codeStr = string(codeFileBytes)
	}

	compiler, err := matchCodeLanguage(session.CodeLangName, session.CodeFile)
	if err != nil { return nil, err }
	err = compiler.Init(codeStr, session.SessionDir)
	if err != nil {
		return nil, err
	}
	return compiler, err
}

// 编译目标程序
func (session *JudgeSession)compileTargetProgram(judgeResult *commonStructs.JudgeResult) error {
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
	// 获取执行指令
	session.Commands = compiler.GetRunArgs()
	session.Compiler = compiler
	return nil
}

// 编译裁判程序
func (session *JudgeSession)compileJudgerProgram(judgeResult *commonStructs.JudgeResult) error {
	_, err := os.Stat(path.Join(session.ConfigDir, session.JudgeConfig.SpecialJudge.Checker))
	if os.IsNotExist(err) {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = fmt.Sprintf("checker file not exists")
		return fmt.Errorf(judgeResult.SeInfo)
	}

	execuable, err := IsExecutableFile(path.Join(session.ConfigDir, session.JudgeConfig.SpecialJudge.Checker))
	if err != nil {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = fmt.Sprintf("checker file not exists")
		return fmt.Errorf("compile checker error: %s", err)
	}
	// 如果是可执行程序，直接放行
	if execuable {
		return nil
	}

	// 判断文件格式
	compiler, err := matchCodeLanguage("auto", session.JudgeConfig.SpecialJudge.Checker)
	if err != nil { return fmt.Errorf("compile checker error: %s", err) }
	switch compiler.(type) {
	case *provider.GnucCompileProvider, *provider.GnucppCompileProvider, *provider.GolangCompileProvider:
		// 初始化编译程序
		codeFile, err := os.Open(path.Join(session.ConfigDir, session.JudgeConfig.SpecialJudge.Checker))
		if err != nil {
			judgeResult.JudgeResult = JudgeFlagSE
			judgeResult.SeInfo = fmt.Sprintf("special judge checker source file open error:\n%s", err.Error())
			return fmt.Errorf(judgeResult.SeInfo)
		}
		defer codeFile.Close()
		code, err := ioutil.ReadAll(codeFile)
		if err != nil {
			judgeResult.JudgeResult = JudgeFlagSE
			judgeResult.SeInfo = fmt.Sprintf("special judge checker source file read error:\n%s", err.Error())
			return fmt.Errorf(judgeResult.SeInfo)
		}
		err = compiler.Init(string(code), session.SessionDir)
		if err != nil {
			return err
		}
		//log.Println(fmt.Sprintf(
		//	"compile (%s) with %s provider",
		//	path.Join(session.ConfigDir, session.JudgeConfig.SpecialJudge.Checker),
		//	compiler.GetName(),
		//))
	default:
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = "special judge checker only support c/c++/go language"
		return fmt.Errorf(judgeResult.SeInfo)
	}

	// 编译程序
	success, ceinfo := compiler.Compile()
	if !success {
		judgeResult.JudgeResult = JudgeFlagSE
		judgeResult.SeInfo = fmt.Sprintf("special judge checker compile error:\n%s", ceinfo)
		return fmt.Errorf("special judge checker compile error:\n%s", ceinfo)
	}
	// 获取执行指令
	session.JudgeConfig.SpecialJudge.Checker = compiler.GetRunArgs()[0]
	session.Compiler = compiler
	return nil
}

