// 生成数据
package packmgr

import (
    "context"
    "fmt"
    "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/executor"
    uuid "github.com/satori/go.uuid"
    "github.com/urfave/cli/v2"
    "io/ioutil"
    "log"
    "os"
    "path"
    "strings"
    "time"
)



func runTestCaseGen(config *structs.JudgeConfiguration, tCase *structs.TestCase, withAnswer bool) error {
    ctx, _ := context.WithTimeout(context.Background(), 3 * time.Second)
    // 如果是generator脚本
    if tCase.UseGenerator {
        inbytes, err := utils.CallGenerator(ctx, tCase, config.ConfigDir)
        if err != nil {
            return err
        }
        // 写入到文件
        err = ioutil.WriteFile(path.Join(config.ConfigDir, tCase.Input), inbytes, 0644)
        if err != nil {
            return err
        }
    }
    return nil
}

// 运行test cases的数据生成
// caseIndex < 0 表示校验全部
func runTestCaseGenerator(config *structs.JudgeConfiguration, caseIndex int, answerCaseIndex uint, withAnswer bool) error {
    // 执行遍历
    if caseIndex < 0 {
        for key, _ := range config.TestCases {
            log.Printf("[generator] run case #%d", key)
            err := runTestCaseGen(config, &config.TestCases[key], withAnswer)
            if err != nil { return err }
        }
    } else {
        log.Printf("[generator] run case #%d", caseIndex)
        err := runTestCaseGen(config, &config.TestCases[caseIndex], withAnswer)
        if err != nil { return err }
    }
    return nil
}

// 运行Testlib的validator校验 (APP入口)
func RunTestCaseGenerator(c *cli.Context) error {
    configFile := c.Args().Get(0)
    _, err := os.Stat(configFile)
    if err != nil && os.IsNotExist(err) {
        return fmt.Errorf("[validator] problem config file (%s) not found", configFile)
    }
    session, err := executor.NewSession(configFile)
    if err != nil { return err }


    if !c.Bool("silence") {
        fmt.Print("Save all the changed to config file? [y/N] ")
        ans := ""
        _, err := fmt.Scanf("%s", &ans)
        if err != nil {
            return nil          // don't crash at EOF
        }
        if len(ans) > 0 && strings.ToLower(ans[:1]) != "y" {
            return nil
        }
    }

    withAnswer := c.Bool("with-answer")
    answerCaseIndex := c.Uint("answer")
    //testCaseIndex := c.Int("case")

    // 编译答案代码
    if withAnswer {
        if len(session.JudgeConfig.AnswerCases) <= 0 {
            return fmt.Errorf("please setup answer case")
        }
        // 强制放到对应地方
        session.SessionId = uuid.NewV4().String()
        session.SessionRoot = "/tmp"
        // 初始化session dir
        sessionDir, err := utils.GetSessionDir(session.SessionRoot, session.SessionId)
        if err != nil {
            return err
        }
        session.SessionDir = sessionDir
        defer session.Clean()
        acase := session.JudgeConfig.AnswerCases[answerCaseIndex]
        // 如果有代码文件
        if acase.FileName != "" {
            session.CodeFile = path.Join(session.ConfigDir, acase.FileName)
        }
        session.CodeLangName = acase.Language
        // 获取对应的编译器提供程序
        compiler, err := session.GetCompiler(acase.Content)
        if err != nil {
            return err
        }
        // 编译程序
        success, ceinfo := compiler.Compile()
        if !success {
            return fmt.Errorf("compile error:\n%s", ceinfo)
        }
        // 获取执行指令
        session.Commands = compiler.GetRunArgs()
        session.Compiler = compiler
    }
    fmt.Println(session.Commands)


    return nil
    //return runTestCaseGenerator(&session.JudgeConfig, testCaseIndex, withAnswer)
}