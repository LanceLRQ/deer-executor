package persistence

import (
	"archive/zip"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
)

// ParseJudgeResultPackageFile parse deer-package file and return a *ProblemProjectPackage object
func ParseJudgeResultPackageFile(packFile string, doValidate bool) (*JudgeResultPackage, error) {
	instance := JudgeResultPackage{
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

// GetResult get judge result (json) from package
func (pack *JudgeResultPackage) GetResult() error {
	return pack.walkDeerPackageBody([]uint8{PackageChunkTypeJudgeResult}, func(typeNum uint8, size int64, reader *io.SectionReader) (int64, error) {
		var err error
		if typeNum == PackageChunkTypeJudgeResult {
			confBytes, err := io.ReadAll(reader)
			if err != nil {
				return 0, err
			}
			utils.JSONBytesObject(confBytes, &pack.JudgeResult)
			return int64(len(confBytes)), nil
		}
		return 0, err
	})
}

// UnpackJudgeResult unpack judge result files to workDir
func (pack *JudgeResultPackage) UnpackJudgeResult(workDir string) error {
	err := pack.walkDeerPackageBody([]uint8{PackageChunkTypeJudgeResult, PackageChunkTypeJudgeResultData}, func(typeNum uint8, size int64, reader *io.SectionReader) (int64, error) {
		var err error
		if typeNum == PackageChunkTypeJudgeResult {
			// unpack file
			configFile := path.Join(workDir, "result.json")
			fp, err := os.Create(configFile)
			if err != nil {
				return size, err
			}
			defer fp.Close()
			if _, err = io.Copy(fp, reader); err != nil {
				return size, err
			}
			return size, nil
		} else if typeNum == PackageChunkTypeJudgeResultData {
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
