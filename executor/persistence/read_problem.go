package persistence

import (
	"github.com/pkg/errors"
	"os"
)

// ParsePackageFile 解析题目包
func ParsePackageFile(packFile string, doValidate bool) (*ProblemProjectPackage, error) {
	instance := ProblemProjectPackage{
		DeerPackageBase: DeerPackageBase{
			presistFilePath: packFile,
		},
	}
	instance.IDeerPackage = &instance

	fileStat, err := os.Stat(packFile)
	if os.IsNotExist(err) {
		return nil, errors.Errorf("file (%s) not exists", packFile)
	} else if err != nil {
		return nil, err
	}
	instance.presistFileSize = fileStat.Size()

	err = instance.parseDeerPackage(doValidate)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}
