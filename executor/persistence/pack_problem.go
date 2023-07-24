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

func NewProblemProjectPackage(conf *commonStructs.JudgeConfiguration) *ProblemProjectPackage {
	instance := ProblemProjectPackage{
		DeerPackageBase: DeerPackageBase{
			Version:   2,
			PackageID: uuid.NewV4(),
		},
		ProblemConfigs: conf,
	}
	instance.IDeerPackage = &instance
	return &instance
}

// WritePackageFile create a problem package file
func (pack *DeerPackageBase) WritePackageFile(presistOptions *ProblemProjectPersisOptions) error {
	pack.presistOptions = presistOptions
	return pack.writePackageFile()
}

// build package body
func (pack *ProblemProjectPackage) buildPackageBody() error {
	options, err := pack.getProblemPersisOptions()
	if err != nil {
		return err
	}
	options.TempFile = fmt.Sprintf("%s.tmp", options.OutFile)
	tempOut, err := os.Create(options.TempFile)
	if err != nil {
		return errors.Errorf("create problem package body file error: %s", err.Error())
	}
	defer tempOut.Close()
	tempOutWriter := bufio.NewWriter(tempOut)
	defer tempOutWriter.Flush()

	// config
	configBytes := utils.ObjectToJSONByte(pack.ProblemConfigs)
	err = pack.writeBodyChunk(tempOutWriter, PackageChunkTypeConfig, uint64(len(configBytes)), bytes.NewBuffer(configBytes))
	if err != nil {
		return err
	}

	// content file
	options.ProblemBodyFile = fmt.Sprintf("%s.zip", options.OutFile)
	err = pack.mergeFilesBinary()
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
func (pack *ProblemProjectPackage) cleanWorkspace() {
	options, err := pack.getProblemPersisOptions()
	if err != nil {
		return
	}
	pack.cleanWorkspaceCommon()
	if options.ProblemBodyFile != "" {
		_ = os.Remove(options.ProblemBodyFile)
	}
	options.ProblemBodyFile = ""
}

// merge problem project files into zip
func (pack *ProblemProjectPackage) mergeFilesBinary() error {
	options, err := pack.getProblemPersisOptions()
	if err != nil {
		return err
	}

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

func (pack *ProblemProjectPackage) getCommonPersisOptions() (*CommonPersisOptions, error) {
	options, err := pack.getProblemPersisOptions()
	if err != nil {
		return nil, err
	}
	return &options.CommonPersisOptions, nil
}

func (pack *ProblemProjectPackage) getProblemPersisOptions() (*ProblemProjectPersisOptions, error) {
	options, ok := pack.presistOptions.(*ProblemProjectPersisOptions)
	if !ok {
		return nil, errors.Errorf("persistOptions must be a *ProblemProjectPersisOptions")
	}
	return options, nil
}
