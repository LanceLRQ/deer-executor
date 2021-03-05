package executor

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/provider"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/pkg/errors"
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
        return provider.NewGnucCompileProvider(), nil
    case "cpp", "gcc-cpp", "gcpp", "g++":
        return provider.NewGnucppCompileProvider(), nil
    case "java":
        return provider.NewJavaCompileProvider(), nil
    case "py2", "python2":
        return provider.NewPy2CompileProvider(), nil
    case "py", "py3", "python3":
        return provider.NewPy3CompileProvider(), nil
    case "php":
        return provider.NewPHPCompileProvider(), nil
    case "go", "golang":
        return provider.NewGolangCompileProvider(), nil
    case "node", "nodejs":
        return provider.NewNodeJSCompileProvider(), nil
    case "rb", "ruby":
        return provider.NewRubyCompileProvider(), nil
    case "auto", "":
        keyword = strings.Replace(path.Ext(fileName), ".", "", -1)
        goto _match
    }
    return nil, errors.Errorf("unsupported language")
}

// 编译文件
// 如果不设置codeStr，默认会读取配置文件里的code_file字段并打开对应文件
func (session *JudgeSession) GetCompiler(codeStr string) (provider.CodeCompileProviderInterface, error) {
    if codeStr == "" {
        codeFileBytes, err := ioutil.ReadFile(session.CodeFile)
        if err != nil {
            return nil, err
        }
        codeStr = string(codeFileBytes)
    }

    compiler, err := matchCodeLanguage(session.CodeLangName, session.CodeFile)
    if err != nil {
        return nil, err
    }
    err = compiler.Init(codeStr, session.SessionDir)
    if err != nil {
        return nil, err
    }
    return compiler, err
}

// 编译目标程序
func (session *JudgeSession) compileTargetProgram(judgeResult *commonStructs.JudgeResult) error {
    // 获取对应的编译器提供程序
    compiler, err := session.GetCompiler("")
    if err != nil {
        judgeResult.JudgeResult = constants.JudgeFlagSE
        judgeResult.SeInfo = err.Error()
        session.Logger.Error(err.Error())
        return err
    }

    // 编译程序
    session.Logger.Infof("Do complie or syntax checkup, Language: %s",  session.CodeLangName)
    success, ceinfo := compiler.Compile()
    if !success {
        judgeResult.JudgeResult = constants.JudgeFlagCE
        judgeResult.CeInfo = ceinfo
        err = errors.Errorf("compile error:\n%s", ceinfo)
        session.Logger.Error(err.Error())
        return err
    }

    // 获取执行指令
    session.Commands = compiler.GetRunArgs()
    session.Compiler = compiler
    return nil
}


// 编译裁判程序
// 如果有已经编译好的裁判程序，则直接返回这个程序
// 打包的时候不会打包二进制文件，重新编译一次
func (session *JudgeSession) compileJudgerProgram(judgeResult *commonStructs.JudgeResult) error {
    // 检查是否存在已经编译好的裁判程序
    cType := "checker"
    if session.JudgeConfig.SpecialJudge.Mode == 2 {
        cType = "interactor"
    }
    cPath, err := utils.GetCompiledBinaryFileAbsPath(cType, session.JudgeConfig.SpecialJudge.Name, session.ConfigDir)
    // 如果有已经编译好的裁判程序，则直接返回这个程序
    if err == nil {
        if s, err := os.Stat(cPath); err == nil && !s.IsDir() {
            session.JudgeConfig.SpecialJudge.Checker = cPath
            return nil
        }
    }

    // 如果没有，则检查checker是否被设置
    jCodeOrExec := path.Join(session.ConfigDir, session.JudgeConfig.SpecialJudge.Checker)
    s, err := os.Stat(jCodeOrExec)
    if os.IsNotExist(err) || s.IsDir() {
        judgeResult.JudgeResult = constants.JudgeFlagSE
        judgeResult.SeInfo = fmt.Sprintf("checker file not exists")
        session.Logger.Error("checker file not exists")
        return errors.Errorf(judgeResult.SeInfo)
    }

    yes, err := utils.IsExecutableFile(jCodeOrExec)
    if err != nil {
        judgeResult.JudgeResult = constants.JudgeFlagSE
        judgeResult.SeInfo = fmt.Sprintf("read checker file error")
        session.Logger.Error(err.Error())
        return err
    } else if yes { // 如果是可执行程序，直接执行
        return nil
    }


    // 编译特判程序
    config := session.JudgeConfig
    binRoot, err := GetOrCreateBinaryRoot(&config)
    if err != nil {
        judgeResult.JudgeResult = constants.JudgeFlagSE
        judgeResult.SeInfo = fmt.Sprintf("create checker bin root error")
        session.Logger.Error(err.Error())
        return err
    }
    session.Logger.Infof("Complie special judge checker, Language: %s",  config.SpecialJudge.CheckerLang)
    compileTarget, err := CompileSpecialJudgeCodeFile(
        config.SpecialJudge.Checker,
        config.SpecialJudge.Name,
        binRoot,
        config.ConfigDir,
        session.LibraryDir,
        config.SpecialJudge.CheckerLang,
    )
    if err != nil {
        judgeResult.JudgeResult = constants.JudgeFlagSE
        judgeResult.SeInfo = fmt.Sprintf("compile checker file error")
        session.Logger.Error(err.Error())
        return err
    }
    // 获取执行指令
    session.JudgeConfig.SpecialJudge.Checker = compileTarget
    return nil
}
