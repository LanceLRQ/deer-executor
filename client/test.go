package client

import (
    "archive/zip"
    "bufio"
    "fmt"
    "github.com/LanceLRQ/deer-common/persistence/problems"
    "github.com/urfave/cli/v2"
)

func Test(c *cli.Context) error {

    zipReader, err := zip.OpenReader("./test.zip")
    if err != nil {
        return err
    }
    defer zipReader.Close()

    file, _, err := problems.FindInZip(zipReader, ".sign")
    if err != nil {
        return err
    }
    bytess := make([]byte, 100)
    buf := bufio.NewReader(*file)
    n, err := buf.Read(bytess)
    if err != nil {
        return err
    }
    fmt.Println(n, bytess)
    return nil
}
