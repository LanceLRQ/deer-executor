package packmgr

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/provider"
    "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-executor/executor"
    "github.com/urfave/cli/v2"
    "os"
    "path"
    "path/filepath"
)

// 针对Testlib支持的编译方法
func compileTestlibCodeFile (source, name, binRoot, configDir, libraryDir, typeName string) error {
    fmt.Printf("build %s [%s]...", typeName, name)
    prefix, ok := constants.TestlibBinaryPrefixs[typeName]
    if !ok { prefix = "" }
    genCodeFile := path.Join(configDir, source)
    compileTarget := path.Join(binRoot, prefix + name)
    _, err := os.Stat(genCodeFile)
    if err != nil && os.IsNotExist(err) {
        return fmt.Errorf("cannot find %s's source code", typeName)
    }
    compiler := provider.NewGnucppCompileProvider()
    ok, ceinfo := compiler.ManualCompile(genCodeFile, compileTarget, [] string{libraryDir})
    if ok {
        fmt.Println("Done.")
    } else {
        fmt.Printf("Error.\n\n%s", ceinfo)
    }
    return nil
}

// 普通的编译方法
func compileNormalCodeFile (source, name, binRoot, configDir, libraryDir, lang, typeName string) error {
    fmt.Printf("build %s [%s]...", typeName, name)
    genCodeFile := path.Join(configDir, source)
    compileTarget := path.Join(binRoot, name)
    _, err := os.Stat(genCodeFile)
    if err != nil && os.IsNotExist(err) {
        return fmt.Errorf("cannot find source code")
    }
    var ok bool
    var ceinfo string
    switch lang {
        case "c", "gcc":
            compiler := provider.NewGnucppCompileProvider()
            ok, ceinfo = compiler.ManualCompile(genCodeFile, compileTarget, [] string{libraryDir})
        case "go", "golang":
            compiler := provider.NewGolangCompileProvider()
            ok, ceinfo = compiler.ManualCompile(genCodeFile, compileTarget)
        default:
            compiler := provider.NewGnucppCompileProvider()
            ok, ceinfo = compiler.ManualCompile(genCodeFile, compileTarget, [] string{libraryDir})
    }
    if ok {
        fmt.Println("Done.")
    } else {
        fmt.Printf("Error.\n\n%s", ceinfo)
    }
    return nil
}

// 编译作业代码
func compileWorkCodeFiles(config structs.JudgeConfiguration, libraryDir string) error {
    binRoot := path.Join(config.ConfigDir, "bin")
    _, err := os.Stat(binRoot)
    if err != nil && os.IsNotExist(err) {
        err = os.MkdirAll(binRoot, 0775)
        if err != nil {
            return fmt.Errorf("cannot create binary work directory: %s", err.Error())
        }
    }
    // generators
    for _, gen := range config.TestLib.Generators {
        err = compileTestlibCodeFile(gen.Source, gen.Name, binRoot, config.ConfigDir, libraryDir, "generator")
        if err != nil {
            return err
        }
    }
    // Validator
    if config.TestLib.Validator != "" && config.TestLib.ValidatorName != "" {
        err = compileTestlibCodeFile(config.TestLib.Validator, config.TestLib.ValidatorName, binRoot, config.ConfigDir, libraryDir, "validator")
        if err != nil {
            return err
        }
    }
    // Checker
    if config.SpecialJudge.Mode > 0 {
        if config.SpecialJudge.Name == "" {
            return fmt.Errorf("please setup special judge checker name")
        }
        if config.SpecialJudge.Checker == "" {
            return fmt.Errorf("please setup special judge checker")
        }
        checkerType := "checker"
        if config.SpecialJudge.Mode == 2 {
            checkerType = "interactor"
        }
        if config.SpecialJudge.UseTestlib {
            err = compileTestlibCodeFile(
                config.SpecialJudge.Checker,
                config.SpecialJudge.Name,
                binRoot,
                config.ConfigDir,
                libraryDir,
                checkerType,
            )
            if err != nil {
                return err
            }
        } else {
            err = compileNormalCodeFile(
                config.SpecialJudge.Checker,
                config.SpecialJudge.Name,
                binRoot,
                config.ConfigDir,
                libraryDir,
                config.SpecialJudge.CheckerLang,
                "special judge " + checkerType,
            )
            if err != nil {
                return err
            }
        }
    }
    return nil
}

// 编译作业代码(APP入口)
func CompileProblemWorkDirSourceCodes(context *cli.Context) error {
    configFile := context.Args().Get(0)
    _, err := os.Stat(configFile)
    if err != nil && os.IsNotExist(err) {
        return fmt.Errorf("problem config file (%s) not found", configFile)
    }
    session, err := executor.NewSession(configFile)
    if err != nil { return err }
    libDir, err := filepath.Abs(context.String("library"))
    if err != nil {
        return err
    }
    err = compileWorkCodeFiles(session.JudgeConfig, libDir)
    return err
}
