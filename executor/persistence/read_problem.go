package persistence

import (
	"archive/zip"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
)

// ParseProblemPackageFile parse deer-package file and return a *ProblemProjectPackage object
func ParseProblemPackageFile(packFile string, doValidate bool) (*ProblemProjectPackage, error) {
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
func (pack *ProblemProjectPackage) GetProblemConfig() error {
	return pack.walkDeerPackageBody([]uint8{PackageChunkTypeConfig}, func(typeNum uint8, size int64, reader *io.SectionReader) (int64, error) {
		var err error
		if typeNum == PackageChunkTypeConfig {
			confBytes, err := io.ReadAll(reader)
			if err != nil {
				return 0, err
			}
			utils.JSONBytesObject(confBytes, &pack.ProblemConfigs)
			return int64(len(confBytes)), nil
		}
		return 0, err
	})
}

// UnpackProblemProject unpack problem config and test files
func (pack *ProblemProjectPackage) UnpackProblemProject(workDir string) error {
	err := pack.walkDeerPackageBody([]uint8{PackageChunkTypeConfig, PackageChunkTypeProject}, func(typeNum uint8, size int64, reader *io.SectionReader) (int64, error) {
		var err error
		if typeNum == PackageChunkTypeConfig {
			// unpack file
			configFile := path.Join(workDir, "problem.json")
			fp, err := os.Create(configFile)
			if err != nil {
				return size, err
			}
			defer fp.Close()
			if _, err = io.Copy(fp, reader); err != nil {
				return size, err
			}
			return size, nil
		} else if typeNum == PackageChunkTypeProject {
			zipReader, err := zip.NewReader(reader, size)
			if err != nil {
				return 0, errors.Errorf("open problem prohect content error: %s", err.Error())
			}
			err = UnZipReader(zipReader, workDir)
			if err != nil {
				return 0, err
			}
		}
		return 0, err
	})
	return err
}
