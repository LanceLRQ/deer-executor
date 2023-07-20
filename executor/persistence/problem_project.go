package persistence

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	commonStructs "github.com/LanceLRQ/deer-executor/v3/executor/structs"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ProblemProjectPersisOptions for persis
type ProblemProjectPersisOptions struct {
	CommonPersisOptions
	ConfigFile      string
	ProjectDir      string
	ProblemBodyFile string
}

// ProblemProjectPackage for problem project
type ProblemProjectPackage struct {
	DeerPackageBase
	// --- body
	problemConfigsBytes []byte // Problem Configs JSON [type: 0x1]
	ProblemBodyTempFile string // Problem package temp file [type: 0x2]
	// --- internal
	ProblemConfigs *commonStructs.JudgeConfiguration
}

func NewProblemProjectPackage(conf *commonStructs.JudgeConfiguration) ProblemProjectPackage {
	return ProblemProjectPackage{
		DeerPackageBase: DeerPackageBase{
			Version:   2,
			PackageID: uuid.NewV4(),
		},
		ProblemConfigs: conf,
	}
}

// WritePackageFile 打包题目数据到文件
func (pack *ProblemProjectPackage) WritePackageFile(options *ProblemProjectPersisOptions) error {
	// Need digital sign?
	if options.DigitalSign {
		if options.DigitalPEM.PublicKey == nil || options.DigitalPEM.PrivateKey == nil {
			return errors.Errorf("digital sign need public key and private key")
		}
	}

	// Build pack body
	err := pack.buildPackageBody(options)
	if err != nil {
		return err
	}

	defer pack.cleanWorkspace(options)

	// Hash body
	err = pack.signPackageBody(&options.CommonPersisOptions)
	if err != nil {
		return err
	}

	// Init outfile
	fout, err := os.Create(options.OutFile)
	if err != nil {
		return errors.Errorf("create problem package file error: %s", err.Error())
	}
	foutWriter := bufio.NewWriter(fout)
	defer fout.Close()

	// Make GPG sign
	if options.DigitalSign {
		err = pack.makeGPGSignature(&options.CommonPersisOptions, foutWriter)
		if err != nil {
			return errors.Errorf("sign problem package file error: %s", err.Error())
		}
	}

	// Write header
	if err = pack.createPackageHeader(foutWriter); err != nil {
		return err
	}

	// Header end with [0 0]
	if _, err := foutWriter.Write([]byte{0x0, 0x0}); err != nil {
		return err
	}

	// Write body
	err = pack.writePackageBody(&options.CommonPersisOptions, foutWriter)
	if err != nil {
		return err
	}

	return nil
}

// build package body
func (pack *ProblemProjectPackage) buildPackageBody(options *ProblemProjectPersisOptions) error {
	options.TempFile = fmt.Sprintf("%s.tmp", options.OutFile)
	tempOut, err := os.Create(options.TempFile)
	if err != nil {
		return errors.Errorf("create problem package body file error: %s", err.Error())
	}
	defer tempOut.Close()
	tempOutWriter := bufio.NewWriter(tempOut)

	// config
	configBytes := utils.ObjectToJSONByte(pack.ProblemConfigs)
	err = pack.writeBodyChunk(tempOutWriter, PackageChunkTypeConfig, uint64(len(configBytes)), bytes.NewBuffer(configBytes))
	if err != nil {
		return err
	}

	// content file
	options.ProblemBodyFile = fmt.Sprintf("%s.zip", options.OutFile)
	err = pack.mergeFilesBinary(options)
	if err != nil {
		return err
	}

	fBody, err := os.Open(options.ProblemBodyFile)
	if err != nil {
		return err
	}
	bodyInfo, err := fBody.Stat()
	if err != nil {
		return err
	}
	err = pack.writeBodyChunk(tempOutWriter, PackageChunkTypeProject, uint64(bodyInfo.Size()), fBody)
	if err != nil {
		return err
	}

	return nil
}

// clean workspace
func (pack *ProblemProjectPackage) cleanWorkspace(options *ProblemProjectPersisOptions) {
	pack.cleanWorkspaceCommon(&options.CommonPersisOptions)
	if options.ProblemBodyFile != "" {
		_ = os.Remove(options.ProblemBodyFile)
	}
	options.ProblemBodyFile = ""
}

// merge problem project files into zip
func (pack *ProblemProjectPackage) mergeFilesBinary(options *ProblemProjectPersisOptions) error {
	zipFile, err := os.Create(options.ProblemBodyFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(options.ProjectDir, func(zpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过配置文件
		if zpath == options.ConfigFile {
			return nil
		}
		// 跳过配置文件夹
		if zpath == options.ProjectDir {
			return nil
		}
		// 跳过二进制文件夹
		binDir := path.Join(options.ProjectDir, "bin")
		if strings.HasPrefix(zpath, binDir) {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(zpath, options.ProjectDir+"/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(zpath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return err
	}
	return nil
}
