/* Deer executor
 * (C) 2019-Now LanceLRQ
 */
package executor

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/provider"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "io/ioutil"
    "os"
    "path"
    "syscall"
)

// 定义ITimer的常量，命名规则遵循Linux的原始设定
const (
    ITimerReal    = 0
    ITimerVirtual = 1
    ITimerVProf   = 2
)

// 定义公共环境变量
var CommonEnvs = []string{"PYTHONIOENCODING=utf-8"}

type ITimerVal struct {
    ItInterval TimeVal
    ItValue    TimeVal
}

type TimeVal struct {
    TvSec  uint64
    TvUsec uint64
}

// 打开文件并获取描述符 (强制文件检查)
func OpenFile(filePath string, flag int, perm os.FileMode) (*os.File, error) {
    if _, err := os.Stat(filePath); err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("file (%s) not exists", filePath)
        } else {
            return nil, fmt.Errorf("open file (%s) error: %s", filePath, err.Error())
        }
    } else {
        if fp, err := os.OpenFile(filePath, flag, perm); err != nil {
            return nil, fmt.Errorf("open file (%s) error: %s", filePath, err.Error())
        } else {
            return fp, nil
        }
    }
}

func Max(x, y int64) int64 {
    if x > y {
        return x
    }
    return y
}

func Max32(a, b int) int {
    if a > b {
        return a
    } else {
        return b
    }
}


// 文件读写(有重试次数，checker专用)
func readFileWithTry(filePath string, name string, tryOnFailed int) ([]byte, string, error) {
    errCnt, errText := 0, ""
    var err error
    for errCnt < tryOnFailed {
        fp, err := OpenFile(filePath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
        if err != nil {
            errText = err.Error()
            errCnt++
            continue
        }
        data, err := ioutil.ReadAll(fp)
        if err != nil {
            _ = fp.Close()
            errText = fmt.Sprintf("Read file(%s) i/o error: %s", name, err.Error())
            errCnt++
            continue
        }
        _ = fp.Close()
        return data, errText, nil
    }
    return nil, errText, err
}

// 检查配置文件里的所有文件是否存在
func CheckRequireFilesExists(config *commonStructs.JudgeConfiguration, configDir string) error {
    var err error
    // 检查特判程序是否存在
    if config.SpecialJudge.Mode != 0 {
        _, err = os.Stat(path.Join(configDir, config.SpecialJudge.Checker))
        if os.IsNotExist(err) {
            return fmt.Errorf("special judge checker file (%s) not exists",config.SpecialJudge.Checker)
        }
    }
    // 检查每个测试数据里的文件是否存在
    // 新版判题机要求无论有没有数据，都要有对应的输入输出文件。
    // 但Testlib模式例外，因为数据是由generator自动生成的。
    for i := 0; i < len(config.TestCases); i++ {
        tcase := config.TestCases[i]
        if !tcase.Enabled || tcase.UseGenerator { continue }
        _, err = os.Stat(path.Join(configDir, tcase.Input))
        if os.IsNotExist(err) {
            return fmt.Errorf("test case (%s) input file (%s) not exists", tcase.Handle, tcase.Input)
        }
        _, err = os.Stat(path.Join(configDir, tcase.Output))
        if os.IsNotExist(err) {
            return fmt.Errorf("test case (%s) output file (%s) not exists", tcase.Handle, tcase.Output)
        }
    }
    return nil
}

// 获取二进制文件的目录
func GetOrCreateBinaryRoot(config *commonStructs.JudgeConfiguration) (string, error) {
    binRoot := path.Join(config.ConfigDir, "bin")
    _, err := os.Stat(binRoot)
    if err != nil && os.IsNotExist(err) {
        err = os.MkdirAll(binRoot, 0775)
        if err != nil {
            return "", fmt.Errorf("cannot create binary work directory: %s", err.Error())
        }
    }
    return binRoot, nil
}

// 普通特殊评测的编译方法
func CompileSpecialJudgeCodeFile (source, name, binRoot, configDir, libraryDir, lang string) (string, error) {
    genCodeFile := path.Join(configDir, source)
    compileTarget := path.Join(binRoot, name)
    _, err := os.Stat(genCodeFile)
    if err != nil && os.IsNotExist(err) {
        return compileTarget, fmt.Errorf("checker source code file not exists")
    }
    var ok bool
    var ceinfo string
    switch lang {
    case "c", "gcc", "gnu-c":
        compiler := provider.NewGnucppCompileProvider()
        ok, ceinfo = compiler.ManualCompile(genCodeFile, compileTarget, [] string{libraryDir})
    case "go", "golang":
        compiler := provider.NewGolangCompileProvider()
        ok, ceinfo = compiler.ManualCompile(genCodeFile, compileTarget)
    case "cpp", "gcc-cpp", "gcpp", "g++", "":
        compiler := provider.NewGnucppCompileProvider()
        ok, ceinfo = compiler.ManualCompile(genCodeFile, compileTarget, [] string{libraryDir})
    default:
        return compileTarget, fmt.Errorf("checker must be written by c/c++/golang")
    }
    if ok {
        return compileTarget, nil
    } else {
        return compileTarget, fmt.Errorf("compile error: %s", ceinfo)
    }
}