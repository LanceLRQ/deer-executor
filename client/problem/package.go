package problem

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/persistence"
    "github.com/LanceLRQ/deer-common/persistence/problems"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/executor"
    "github.com/urfave/cli/v2"
    "log"
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
            s, _, err := problems.ReadProblemInfo(configFile, false, "")
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
