package persistence

import (
	"encoding/binary"
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path"
)

/*********
------------------------
|MAG|VER|RSZ|BSZ|SSZ|CMP| Signature | Result | Body
------------------------
| 2 | 1 | 4 | 4 | 2 | 1 | ...
------------------------
**********/
type JudgeResultPackage struct {
	Version 		uint8				// (VER) Package Version
	ResultSize 		uint32				// (RSZ) Result JSON Text Size
	BodySize 		uint32				// (BSZ) Result Body Size
	SignSize		uint16				// (SSZ) Signature Size
	CertSize		uint16				// (CSZ) Public Certificate Size
	Compressor		uint8				// (CMP) Compressor type: 0-disabled; 1-zip
	Certificate		[]byte				// Public Certificate
	Signature       []byte		 		// Signature: SHA256(Result + Body)
	Result 			[]byte				// Result JSON
	Body			[]byte				// Body Binary
}

type JudgeResultPackageBody struct {
	Size 			uint32
	FileName 		string
	Content			[]byte
}


func readAndWriteToTempFile(writer io.Writer, filePath string) error {
	buf32 := make([]byte, 4)
	body, err := executor.ReadFile(filePath)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(buf32, uint32(len(body)))
	if _, err := writer.Write(buf32); err != nil {
		return fmt.Errorf("write temp file error: %s", err.Error())
	}
	if _, err := writer.Write([]byte(filePath + "\n")); err != nil {
		return fmt.Errorf("write temp file error: %s", err.Error())
	}
	if _, err := writer.Write(body); err != nil {
		return fmt.Errorf("write temp file error: %s", err.Error())
	}
	return nil
}

func MergeResultBinary(judgeResult *executor.JudgeResult, compressType uint8) (string, error) {
	tmpFileName := uuid.NewV1().String() + ".tmp"
	tmpFilePath := path.Join("/tmp/", tmpFileName)
	tmpFile, err := os.OpenFile(tmpFilePath, os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("create temp file error: %s", err.Error())
	}
	//defer (func() {
	//	_ = os.Remove(tmpFilePath)
	//})()

	for _, testCase := range judgeResult.TestCases {
		err = readAndWriteToTempFile(tmpFile, testCase.ProgramOut)
		_ = readAndWriteToTempFile(tmpFile, testCase.ProgramError)
		_ = readAndWriteToTempFile(tmpFile, testCase.ProgramLog)
		_ = readAndWriteToTempFile(tmpFile, testCase.JudgerOut)
		_ = readAndWriteToTempFile(tmpFile, testCase.JudgerError)
		_ = readAndWriteToTempFile(tmpFile, testCase.JudgerLog)
		_ = readAndWriteToTempFile(tmpFile, testCase.JudgerReport)
	}
	_ = tmpFile.Close()

	// TODO Zip
	if compressType == 1 {

	}

	return tmpFilePath, nil
}

//func PersistentJudgeResult(
//	session *executor.JudgeSession,
//	judgeResult *executor.JudgeResult,
//	certFile string,
//	certKeyFile string,
//	compress bool,
//	outFile string,
//) error {
//	relPackage := JudgeResultPackage{}
//	relPackage.Version = 1
//	relPackage.Result = executor.ObjectToJSONByte(judgeResult)
//
//
//}