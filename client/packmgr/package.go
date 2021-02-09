package packmgr

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/persistence"
    "github.com/LanceLRQ/deer-common/persistence/problems"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/v2/executor"
    "github.com/pkg/errors"
    "github.com/urfave/cli/v2"
    "log"
    "os"
    "path"
    "strings"
)

func BuildProblemPackage(c *cli.Context) error {

    if c.String("passphrase") != "" {
        log.Println("[warn] Using a password on the command line interface can be insecure.")
    }
    passphrase := []byte(c.String("passphrase"))
    configFile := c.Args().Get(0)
    outputFile := c.Args().Get(1)

    if c.Bool("zip") && !strings.HasSuffix(configFile, "problem.json") {
        return errors.Errorf("config file must named 'problem.json' in zip mode")
    }

    var err error
    var pem *persistence.DigitalSignPEM

    _, err = os.Stat(configFile)
    if err != nil && os.IsNotExist(err) {
        return errors.Errorf("problem config file (%s) not found", configFile)
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

    if c.Bool("zip") {
        err = problems.PackProblemsAsZip(&options)
        if err != nil {
            return err
        }
    } else {
        err = problems.PackProblems(&session.JudgeConfig, &options)
        if err != nil {
            return err
        }
    }
    return nil
}

// 题目包解包
func UnpackProblemPackage(c *cli.Context) error {
    packageFile := c.Args().Get(0)
    workDir := c.Args().Get(1)
    // 如果路径存在目录或者文件
    if _, err := os.Stat(workDir); err == nil {
        return errors.Errorf("work directory (%s) path exisis", workDir)
    }
    // 创建目录
    if err := os.MkdirAll(workDir, 0775); err != nil {
        return err
    }
    if c.Bool("no-validate") {
        log.Println("[warn] package validation had been disabled!")
    }
    // 检查是否为题目包
    isDeerPack, err := utils.IsProblemPackage(packageFile)
    if err != nil {
        return err
    }
    isZip, err := utils.IsZipFile(packageFile)
    if err != nil {
        return err
    }
    // 解包
    if isDeerPack {
        if _, _, err := problems.ReadProblemInfo(packageFile, true, !c.Bool("no-validate"), workDir); err != nil {
            return err
        }
    } else if isZip {
        if _, _, err := problems.ReadProblemInfoZip(packageFile, true, !c.Bool("no-validate"), workDir); err != nil {
            return err
        }
        // clean meta file
        _ = os.Remove(path.Join(workDir, ".sign"))
        _ = os.Remove(path.Join(workDir, ".gpg"))
    } else {
        return errors.Errorf("not a deer-executor problem package file")
    }
    fmt.Println("Done.")
    return nil
}

// 访问题目包信息
func ReadProblemInfo(c *cli.Context) error {
    packageFile := c.Args().Get(0)
    isDeerPack, err := utils.IsProblemPackage(packageFile)
    if err != nil {
        return err
    }
    isZip, err := utils.IsZipFile(packageFile)
    if err != nil {
        return err
    }
    // 如果是题目包文件，进行解包
    if isDeerPack {
        if c.Bool("gpg") {
            g, err := problems.ReadProblemGPGInfo(packageFile)
            if err != nil {
                fmt.Println(err.Error())
                return nil
            }
            fmt.Println(g)
        } else {
            s, _, err := problems.ReadProblemInfo(packageFile, false, false, "")
            if err != nil {
                return err
            }
            fmt.Println(utils.ObjectToJSONStringFormatted(s))
        }
    } else if isZip {
        if c.Bool("gpg") {
            g, err := problems.ReadProblemGPGInfoZip(packageFile)
            if err != nil {
                fmt.Println(err.Error())
                return nil
            }
            fmt.Println(g)
        } else {
            s, _, err := problems.ReadProblemInfoZip(packageFile, false, false, "")
            if err != nil {
                return err
            }
            fmt.Println(utils.ObjectToJSONStringFormatted(s))
        }
    } else {
        return errors.Errorf("not a deer-executor problem package file")
    }

    return nil
}
