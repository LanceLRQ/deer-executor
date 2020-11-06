package packmgr

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/persistence"
    "github.com/LanceLRQ/deer-common/persistence/problems"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/executor"
    "github.com/urfave/cli/v2"
    "log"
    "os"
)

func BuildProblemPackage(c *cli.Context) error {

    if c.String("passphrase") != "" {
        log.Println("[warn] Using a password on the command line interface can be insecure.")
    }
    passphrase := []byte(c.String("passphrase"))
    configFile := c.Args().Get(0)
    outputFile := c.Args().Get(1)

    var err error
    var pem *persistence.DigitalSignPEM

    _, err = os.Stat(configFile)
    if err != nil && os.IsNotExist(err) {
        return fmt.Errorf("problem config file (%s) not found", configFile)
    }

    if c.Bool("sign") {
        pem, err = persistence.GetArmorPublicKey(c.String("gpg-key"), passphrase)
        if err != nil {
            return err
        }
    }
    options := persistence.ProblemPackageOptions{}
    options.ConfigFile = configFile
    options.DigitalSign = c.Bool("sign")
    options.DigitalPEM = pem
    options.OutFile = outputFile

    // problem
    session, err := executor.NewSession(configFile)
    if err != nil {
        return err
    }
    options.ConfigDir = session.ConfigDir

    err = executor.CheckRequireFilesExists(&session.JudgeConfig, options.ConfigDir)
    if err != nil {
        return err
    }

    err = problems.PackProblems(&session.JudgeConfig, &options)
    if err != nil {
        return err
    }

    return nil
}

// 题目包解包
func UnpackProblemPackage(c *cli.Context) error {
    packageFile := c.Args().Get(0)
    workDir := c.Args().Get(1)
    // 如果路径存在目录或者文件
    if _, err := os.Stat(workDir); err == nil {
        return fmt.Errorf("work directory (%s) path exisis", workDir)
    }
    // 检查题目包是否存在
    yes, err := problems.IsProblemPackage(packageFile)
    if err != nil {
        return err
    }
    if !yes {
        return fmt.Errorf("not a problem package")
    }
    // 创建目录
    if err := os.MkdirAll(workDir, 0775); err != nil {
        return err
    }
    if c.Bool("no-validate") {
        log.Println("[warn] package validation had been disabled!")
    }
    // 解包
    if _, _, err := problems.ReadProblemInfo(packageFile, true, !c.Bool("no-validate"), workDir); err != nil {
        return err
    }
    fmt.Println("Done.")
    return nil
}

// 访问题目包信息
func ReadProblemInfo(c *cli.Context) error {
    configFile := c.Args().Get(0)
    yes, err := problems.IsProblemPackage(configFile)
    if err != nil {
        return err
    }
    // 如果是题目包文件，进行解包
    if yes {
        if c.Bool("sign") {
            g, err := problems.ReadProblemGPGInfo(configFile)
            if err != nil {
                return err
            }
            fmt.Println(g)
        } else {
            s, _, err := problems.ReadProblemInfo(configFile, false, true, "")
            if err != nil {
                return err
            }
            fmt.Println(utils.ObjectToJSONStringFormatted(s))
        }
    } else {
        return fmt.Errorf("not deer-executor problem package file")
    }

    return nil
}
