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


// 编译作业代码
func compileWorkCodeFiles(config structs.JudgeConfiguration, libraryDir string) error {
    binRoot, err := executor.GetOrCreateBinaryRoot(&config)
    if err != nil {
        return err
    }
    // Generators
    if config.TestLib.Generators != nil {
        for _, gen := range config.TestLib.Generators {
            err = compileTestlibCodeFile(gen.Source, gen.Name, binRoot, config.ConfigDir, libraryDir, "generator")
            if err != nil {
                return err
            }
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
            fmt.Printf("build %s [%s]...", "special judge " + checkerType, config.SpecialJudge.Name)
            _, err = executor.CompileSpecialJudgeCodeFile(
                config.SpecialJudge.Checker,
                config.SpecialJudge.Name,
                binRoot,
                config.ConfigDir,
                libraryDir,
                config.SpecialJudge.CheckerLang,
            )
            if err != nil {
                fmt.Printf("Error!\n%s", err.Error())
                return fmt.Errorf("compile error")
            } else {
                fmt.Println("Ok!")
            }
        }
    }
    return nil
}

// 编译作业代码(APP入口)
func CompileProblemWorkDirSourceCodes(c *cli.Context) error {
    configFile := c.Args().Get(0)
    _, err := os.Stat(configFile)
    if err != nil && os.IsNotExist(err) {
        return fmt.Errorf("problem config file (%s) not found", configFile)
    }
    session, err := executor.NewSession(configFile)
    if err != nil { return err }
    libDir, err := filepath.Abs(c.String("library"))
    if err != nil {
        return fmt.Errorf("get library root error: %s", err.Error())
    }
    if s, err := os.Stat(libDir); err != nil {
        return fmt.Errorf("library root not exists")
    } else {
        if !s.IsDir() {
            return fmt.Errorf("library root not a directory")
        }
    }
    err = compileWorkCodeFiles(session.JudgeConfig, libDir)
    return err
}
