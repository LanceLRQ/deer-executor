package client

import (
    "fmt"
    "os"
    "path"
)

// create and get session directory
func GetSessionDir(workDir string, sessionId string) (string, error) {
    _, err := os.Stat(workDir)
    if os.IsNotExist(err) {
        return "", fmt.Errorf("work dir (%s) not exists", workDir)
    } else if err != nil {
        return "", err
    }
    sessionDir := path.Join(workDir, sessionId)
    _, err = os.Stat(sessionDir)
    if os.IsExist(err) {
        _ = os.RemoveAll(sessionDir)
    }
    err = os.Mkdir(sessionDir, 0755)
    if err != nil {
        return "", fmt.Errorf("create session dir error: %s", err.Error())
    }
    return sessionDir, nil
}
