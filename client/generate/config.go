package generate

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/provider"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/urfave/cli/v2"
    "log"
    "os"
)

// 生成编译器配置(程序使用)
func MakeCompileConfigFile(c *cli.Context) error {
    config := provider.CompileCommands
    output := c.String("output")
    if output == "" {
        output = "./compilers.json"
    }
    s, err := os.Stat(output)
    if s != nil || os.IsExist(err) {
        log.Fatal("output file exists")
        return nil
    }
    fp, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        log.Fatalf("open output file error: %s\n", err.Error())
        return nil
    }
    defer fp.Close()
    _, err = fp.WriteString(utils.ObjectToJSONStringFormatted(config))
    if err != nil {
        return err
    }
    return nil
}

// 生成JIT内存宽限配置
func MakeJITMemoryConfigFile(c *cli.Context) error {
    config := constants.MemorySizeForJIT
    output := c.String("output")
    if output == "" {
        output = "./jit_memory.json"
    }
    s, err := os.Stat(output)
    if s != nil || os.IsExist(err) {
        log.Fatal("output file exists")
        return nil
    }
    fmt.Println(output)
    fp, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        log.Fatalf("open output file error: %s\n", err.Error())
        return nil
    }
    defer fp.Close()
    _, err = fp.WriteString(utils.ObjectToJSONStringFormatted(config))
    if err != nil {
        return err
    }
    return nil
}

