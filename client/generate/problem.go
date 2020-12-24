package generate

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/persistence/problems"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/executor"
    "github.com/urfave/cli/v2"
    "io/ioutil"
    "os"
    "path"
)

func makeProblmConfig() (*commonStructs.JudgeConfiguration, error) {
    session, err := executor.NewSession("")
    if err != nil { return nil, err }
    config := session.JudgeConfig
    config.TestCases = []commonStructs.TestCase{{}}
    config.Limitation = make(map[string]commonStructs.JudgeResourceLimit)
    config.Limitation["gcc"] = commonStructs.JudgeResourceLimit{
        TimeLimit:     config.TimeLimit,
        MemoryLimit:   config.MemoryLimit,
        RealTimeLimit: config.RealTimeLimit,
        FileSizeLimit: config.FileSizeLimit,
    }
    config.AnswerCases = []commonStructs.AnswerCase {{},}
    config.SpecialJudge.CheckerCases = []commonStructs.SpecialJudgeCheckerCase{{}}
    config.Problem.Sample = []commonStructs.ProblemIOSample{{}}
    config.TestLib.ValidatorCases = []commonStructs.TestlibValidatorCase{{}}
    config.TestLib.Generators = []commonStructs.TestlibGenerator{{}}
    return &config, nil
}

// 生成评测配置文件
func MakeProblemConfigFile(c *cli.Context) error {
    config, err := makeProblmConfig()
    if err != nil { return err }
    output := c.String("output")
    if output != "" {
        s, err := os.Stat(output)
        if s != nil || os.IsExist(err) {
            return fmt.Errorf("output file exists")
        }
        fp, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
        if err != nil {
            return fmt.Errorf("open output file error: %s\n", err.Error())
        }
        defer fp.Close()
        _, err = fp.WriteString(utils.ObjectToJSONStringFormatted(config))
        if err != nil {
            return err
        }
    } else {
        fmt.Println(utils.ObjectToJSONStringFormatted(config))
    }
    return nil
}

// 创建一个题目工作目录
func InitProblemWorkDir(c *cli.Context) error {
    workDir := c.Args().Get(0)
    // 如果路径存在目录或者文件
    if _, err := os.Stat(workDir); err == nil {
        return fmt.Errorf("work directory (%s) path exisis", workDir)
    }
    // 创建目录
    if err := os.MkdirAll(workDir, 0775); err != nil {
        return err
    }
    example := c.String("name")
    if example != "" {
        packageFile := path.Join("./lib/example", example)
        // 检查题目包是否存在
        yes, err := problems.IsProblemPackage(packageFile)
        if err != nil {
            return err
        }
        if !yes {
            return fmt.Errorf("not a problem package")
        }
        // 如果指定了对应的模板
        if _, _, err := problems.ReadProblemInfo(packageFile, true, true, workDir); err != nil {
            return err
        }
    } else {
        // 创建文件夹
        dirs := []string{"answers", "cases", "bin", "codes", "generators"}
        for _, dirname := range dirs {
            err := os.MkdirAll(path.Join(workDir, dirname), 0775)
            if err != nil {
                return err
            }
        }
        /// 创建配置
        config, err := makeProblmConfig()
        if err != nil {
            return err
        }
        // 写入到文件
        if err = ioutil.WriteFile(path.Join(workDir, "problem.json"), []byte(utils.ObjectToJSONStringFormatted(config)), 0664); err != nil {
            return err
        }
    }
    return nil
}