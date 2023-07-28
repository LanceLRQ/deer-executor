package persistence

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/LanceLRQ/deer-executor/v3/executor/constants"
	commonStructs "github.com/LanceLRQ/deer-executor/v3/executor/structs"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path"
)

func NewJudgeResultPackage(judgeResult *commonStructs.JudgeResult) *JudgeResultPackage {
	instance := JudgeResultPackage{
		DeerPackageBase: DeerPackageBase{
			Version:   1,
			PackageID: uuid.NewV4(),
		},
		JudgeResult: judgeResult,
	}
	instance.IDeerPackage = &instance
	return &instance
}

// WritePackageFile create a judge result package file
func (pack *JudgeResultPackage) WritePackageFile(presistOptions *JudgeResultPersisOptions) error {
	pack.presistOptions = presistOptions
	return pack.writePackageFile()
}

// build package body
func (pack *JudgeResultPackage) buildPackageBody() error {
	options, err := pack.getJudgeResultPersisOptions()
	if err != nil {
		return err
	}
	options.TempFile = fmt.Sprintf("%s.tmp", options.OutFile)
	tempOut, err := os.Create(options.TempFile)
	if err != nil {
		return errors.Errorf("create judge result package body file error: %s", err.Error())
	}
	defer tempOut.Close()

	// config
	configBytes := utils.ObjectToJSONByte(pack.JudgeResult)
	err = pack.writeBodyChunk(tempOut, PackageChunkTypeJudgeResult, uint64(len(configBytes)), bytes.NewBuffer(configBytes))
	if err != nil {
		return err
	}

	// content file
	options.JudgeResultDataFile = fmt.Sprintf("%s.tmp.zip", options.OutFile)
	err = pack.mergeFilesBinary()
	if err != nil {
		return err
	}
	fBody, err := os.Open(options.JudgeResultDataFile)
	if err != nil {
		return err
	}
	bodyInfo, err := fBody.Stat()
	if err != nil {
		return err
	}
	err = pack.writeBodyChunk(tempOut, PackageChunkTypeJudgeResultData, uint64(bodyInfo.Size()), fBody)
	if err != nil {
		return err
	}

	return nil
}

// merge problem project files into zip
func (pack *JudgeResultPackage) mergeFilesBinary() error {
	options, err := pack.getJudgeResultPersisOptions()
	if err != nil {
		return err
	}

	zipFile, err := os.Create(options.JudgeResultDataFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	addFile := func(fileName string) error {
		workDir := options.SessionDir
		fp, err := os.Open(path.Join(workDir, fileName))
		if err != nil {
			return err
		}
		defer fp.Close()
		fileInfo, err := fp.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}
		header.Name = fileName
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, fp)
		return err
	}

	for _, testCase := range pack.JudgeResult.TestCases {
		// if you don't need to save accepted test data
		if !options.SaveAcceptedData && testCase.JudgeResult == constants.JudgeFlagAC {
			continue
		}

		fileList := []string{
			testCase.ProgramOut, testCase.ProgramError,
			testCase.CheckerOut, testCase.CheckerError, testCase.CheckerReport,
		}

		for _, item := range fileList {
			err = addFile(item)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// clean workspace
func (pack *JudgeResultPackage) cleanWorkspace() {
	options, err := pack.getJudgeResultPersisOptions()
	if err != nil {
		return
	}
	pack.cleanWorkspaceCommon()
	//if options.JudgeResultDataFile != "" {
	//	_ = os.Remove(options.JudgeResultDataFile)
	//}
	options.JudgeResultDataFile = ""
}

func (pack *JudgeResultPackage) getCommonPersisOptions() (*CommonPersisOptions, error) {
	options, err := pack.getJudgeResultPersisOptions()
	if err != nil {
		return nil, err
	}
	return &options.CommonPersisOptions, nil
}

func (pack *JudgeResultPackage) getJudgeResultPersisOptions() (*JudgeResultPersisOptions, error) {
	options, ok := pack.presistOptions.(*JudgeResultPersisOptions)
	if !ok {
		return nil, errors.Errorf("persistOptions must be a *JudgeResultPersisOptions")
	}
	return options, nil
}
