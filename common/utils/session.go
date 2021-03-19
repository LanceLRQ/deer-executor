package utils

import (
	"github.com/pkg/errors"
	"os"
	"path"
)

// create and get session directory
func GetSessionDir(workDir string, SessionID string) (string, error) {
	_, err := os.Stat(workDir)
	if os.IsNotExist(err) {
		return "", errors.Errorf("work dir (%s) not exists", workDir)
	} else if err != nil {
		return "", err
	}
	sessionDir := path.Join(workDir, SessionID)
	s, err := os.Stat(sessionDir)
	if s != nil || os.IsExist(err) {
		_ = os.RemoveAll(sessionDir)
	}
	err = os.Mkdir(sessionDir, 0755)
	if err != nil {
		return "", errors.Errorf("create session dir error: %s", err.Error())
	}
	return sessionDir, nil
}
