package judge_result

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/LanceLRQ/deer-executor/persistence"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func readAndWriteToTempFile(writer io.Writer, fileName string, workDir string) error {
	buf16 := make([]byte, 16)
	buf32 := make([]byte, 4)
	body, err := ioutil.ReadFile(path.Join(workDir, fileName))
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint16(buf16, persistence.JudgeBodyPackageMagicCode)
	binary.BigEndian.PutUint32(buf32, uint32(len(body)))
	if _, err := writer.Write(buf16); err != nil {
		return fmt.Errorf("write temp file error: %s", err.Error())
	}
	if _, err := writer.Write(buf32); err != nil {
		return fmt.Errorf("write temp file error: %s", err.Error())
	}
	if _, err := writer.Write([]byte(fileName + "\n")); err != nil {
		return fmt.Errorf("write temp file error: %s", err.Error())
	}
	if _, err := writer.Write(body); err != nil {
		return fmt.Errorf("write temp file error: %s", err.Error())
	}
	return nil
}

func mergeResultBinary(session *executor.JudgeSession, judgeResult *executor.JudgeResult, compressType uint8) (string, error) {
	tmpFileName := uuid.NewV1().String() + ".tmp"
	tmpFilePath := path.Join("/tmp/", tmpFileName)
	var testCaseWriter io.Writer
	tmpFile, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("create temp file error: %s", err.Error())
	}
	defer tmpFile.Close()
	if compressType == 1 { // GZIP
		zipWriter := gzip.NewWriter(tmpFile)
		testCaseWriter = zipWriter
		defer zipWriter.Close()
	} else {
		testCaseWriter = tmpFile
	}

	for _, testCase := range judgeResult.TestCases {
		err = readAndWriteToTempFile(testCaseWriter, testCase.ProgramOut, session.SessionDir)
		_ = readAndWriteToTempFile(testCaseWriter, testCase.ProgramError, session.SessionDir)
		_ = readAndWriteToTempFile(testCaseWriter, testCase.JudgerOut, session.SessionDir)
		_ = readAndWriteToTempFile(testCaseWriter, testCase.JudgerError, session.SessionDir)
		_ = readAndWriteToTempFile(testCaseWriter, testCase.JudgerReport, session.SessionDir)
	}

	return tmpFilePath, nil
}

func writeFileHeaderAndResult (writer io.Writer, pack JudgeResultPackage) error {
	buf8 := make([]byte, 1)
	buf16 := make([]byte, 2)
	buf32 := make([]byte, 4)

	// magic
	binary.BigEndian.PutUint16(buf16, persistence.JudgeResultMagicCode)
	if _, err := writer.Write(buf16); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}

	// Version
	buf8[0] = pack.Version
	if _, err := writer.Write(buf8); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}

	// CompressorType
	buf8[0] = pack.CompressorType
	if _, err := writer.Write(buf8); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}

	// ResultSize
	binary.BigEndian.PutUint32(buf32, pack.ResultSize)
	if _, err := writer.Write(buf32); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}

	// BodySize
	binary.BigEndian.PutUint32(buf32, pack.BodySize)
	if _, err := writer.Write(buf32); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}

	// CertSize
	binary.BigEndian.PutUint16(buf16, pack.CertSize)
	if _, err := writer.Write(buf16); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}

	// Certificate
	if pack.CertSize > 0 {
		if _, err := writer.Write(pack.Certificate); err != nil {
			return fmt.Errorf("write result file error: %s", err.Error())
		}
	}

	return nil
}

func PersistentJudgeResult(
	session *executor.JudgeSession,
	judgeResult *executor.JudgeResult,
	options JudgeResultPersisOptions,
) error {
	fout, err := os.Create(options.OutFile)
	if err != nil {
		return fmt.Errorf("create result file error: %s", err.Error())
	}
	defer fout.Close()

	if options.DigitalSign {
		if options.DigitalPEM.PublicKey == nil || options.DigitalPEM.PrivateKey == nil {
			return fmt.Errorf("digital sign need public key and private key")
		}
	}

	resultBytes := executor.ObjectToJSONByte(judgeResult)

	bodyFile, err := mergeResultBinary(session, judgeResult, options.CompressorType)
	if err != nil {
		return err
	}

	bodyInfo, err := os.Stat(bodyFile)
	if err != nil { return err }

	certSize := 0			// 0 means disable cert
	if options.DigitalSign {
		certSize = len(options.DigitalPEM.PublicKeyRaw)
	}

	pack := JudgeResultPackage{
		Version: 1,
		Result: resultBytes,
		ResultSize: uint32(len(resultBytes)),
		BodySize: uint32(bodyInfo.Size()),
		CertSize: uint16(certSize),
		Certificate: options.DigitalPEM.PublicKeyRaw,
		CompressorType: options.CompressorType,
	}
	// Write Header
	err = writeFileHeaderAndResult(fout, pack)
	if err != nil {
		return err
	}

	// Write Signature
	fBody, err := os.Open(bodyFile)
	if err != nil { return err }

	hash, err := persistence.SHA256Streams([]io.Reader{
		bytes.NewReader(resultBytes),
		fBody,
	})
	if err != nil { return err }
	_ = fBody.Close()
	if options.DigitalSign {
		hash, err = persistence.RSA2048Sign(hash, options.DigitalPEM.PrivateKey)
		if err != nil { return err }
	}
	buf16 := make([]byte, 2)
	signSize := uint16(len(hash))
	// SignSize
	binary.BigEndian.PutUint16(buf16, signSize)
	if _, err := fout.Write(buf16); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}
	// Signature
	if _, err := fout.Write(hash); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}

	// Write Result and Body
	// 要注意先写入result，再写body，方便后续校验的时候直接顺序读取
	if _, err := fout.Write(pack.Result); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}
	fBody, err = os.Open(bodyFile)
	if err != nil { return err }
	defer fBody.Close()
	// Copy Body to fout
	if _, err := io.Copy(fout, fBody); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}


	return nil

}