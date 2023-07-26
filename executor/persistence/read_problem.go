package persistence

import (
	"archive/zip"
	"bytes"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
	"path"
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
func (pack *ProblemProjectPackage) GetProblemConfig() error {
	return pack.walkDeerPackageBody([]interface{}{PackageChunkTypeConfig}, func(typeNum uint8, size int64, reader *io.SectionReader) (int64, error) {
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
}

// GetProblemGPGInfo get GPG certification info
func (pack *ProblemProjectPackage) GetProblemGPGInfo() (string, error) {
	if pack.GPGCertSize == 0 {
		return "no GPG public key", nil
	}
	elist, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(pack.GPGCertificate))
	if err != nil {
		return "", err
	}
	if len(elist) < 1 {
		return "", errors.Errorf("GPG key error")
	}
	rel := ""
	for _, identify := range elist[0].Identities {
		rel += identify.Name + "\n"
	}
	return rel, nil
}

// UnpackProblemProject unpack problem config and test files
func (pack *ProblemProjectPackage) UnpackProblemProject(workDir string) error {
	err := pack.walkDeerPackageBody([]interface{}{PackageChunkTypeConfig, PackageChunkTypeProject}, func(typeNum uint8, size int64, reader *io.SectionReader) (int64, error) {
		var err error
		if typeNum == PackageChunkTypeConfig {
			confBytes, err := io.ReadAll(reader)
			if err != nil {
				return 0, err
			}
			pack.problemConfigsBytes = confBytes
			utils.JSONBytesObject(confBytes, &pack.ProblemConfigs)
			size := int64(len(confBytes))
			// unpack file
			configFile := path.Join(workDir, "problem.json")
			fp, err := os.Create(configFile)
			if err != nil {
				return size, err
			}
			defer fp.Close()
			_, err = fp.Write(pack.problemConfigsBytes)
			if err != nil {
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
