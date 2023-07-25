package persistence

import (
	"github.com/LanceLRQ/deer-executor/v3/executor/structs"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	"io"
	"os"
)

// ParsePackageFile parse deer-package file and return a *ProblemProjectPackage object
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

// GetProblemConfig get problem config from package
func (pack *ProblemProjectPackage) GetProblemConfig() (*structs.JudgeConfiguration, error) {
	err := pack.walkDeerPackageBody([]interface{}{PackageChunkTypeConfig}, func(typeNum uint8, size int64, reader io.Reader) (int64, error) {
		var err error
		if typeNum == PackageChunkTypeConfig {
			confBytes, err := io.ReadAll(reader)
			if err != nil {
				return 0, err
			}
			pack.problemConfigsBytes = confBytes
			utils.JSONBytesObject(confBytes, &pack.ProblemConfigs)
			return int64(len(confBytes)), nil
		}
		return 0, err
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
