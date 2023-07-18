package result

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	constants2 "github.com/LanceLRQ/deer-executor/v3/executor/constants"
	persistence2 "github.com/LanceLRQ/deer-executor/v3/executor/persistence"
	commonStructs "github.com/LanceLRQ/deer-executor/v3/executor/structs"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
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
	binary.BigEndian.PutUint16(buf16, constants2.JudgeBodyPackageMagicCode)
	binary.BigEndian.PutUint32(buf32, uint32(len(body)))
	if _, err := writer.Write(buf16); err != nil {
		return errors.Errorf("write temp file error: %s", err.Error())
	}
	if _, err := writer.Write(buf32); err != nil {
		return errors.Errorf("write temp file error: %s", err.Error())
	}
	if _, err := writer.Write([]byte(fileName + "\n")); err != nil {
		return errors.Errorf("write temp file error: %s", err.Error())
	}
	if _, err := writer.Write(body); err != nil {
		return errors.Errorf("write temp file error: %s", err.Error())
	}
	return nil
}

func mergeResultBinary(
	options *persistence2.JudgeResultPersisOptions,
	judgeResult *commonStructs.JudgeResult,
) (string, error) {
	tmpFileName := uuid.NewV1().String() + ".tmp"
	tmpFilePath := path.Join("/tmp/", tmpFileName)
	var testCaseWriter io.Writer
	tmpFile, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", errors.Errorf("create temp file error: %s", err.Error())
	}
	defer tmpFile.Close()
	if options.CompressorType == 1 { // GZIP
		zipWriter := gzip.NewWriter(tmpFile)
		testCaseWriter = zipWriter
		defer zipWriter.Close()
	} else {
		testCaseWriter = tmpFile
	}

	for _, testCase := range judgeResult.TestCases {
		// 如果不需要保留AC的数据
		if !options.SaveAcceptedData && testCase.JudgeResult == constants2.JudgeFlagAC {
			continue
		}
		err = readAndWriteToTempFile(testCaseWriter, testCase.ProgramOut, options.SessionDir)
		_ = readAndWriteToTempFile(testCaseWriter, testCase.ProgramError, options.SessionDir)
		_ = readAndWriteToTempFile(testCaseWriter, testCase.CheckerOut, options.SessionDir)
		_ = readAndWriteToTempFile(testCaseWriter, testCase.CheckerError, options.SessionDir)
		_ = readAndWriteToTempFile(testCaseWriter, testCase.CheckerReport, options.SessionDir)
	}

	return tmpFilePath, nil
}

func writeFileHeaderAndResult(writer io.Writer, pack JudgeResultPackage) error {
	buf8 := make([]byte, 1)
	buf16 := make([]byte, 2)
	buf32 := make([]byte, 4)

	// magic
	binary.BigEndian.PutUint16(buf16, constants2.JudgeResultMagicCode)
	if _, err := writer.Write(buf16); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}

	// Version
	buf8[0] = pack.Version
	if _, err := writer.Write(buf8); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}

	// CompressorType
	buf8[0] = pack.CompressorType
	if _, err := writer.Write(buf8); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}

	// ResultSize
	binary.BigEndian.PutUint32(buf32, pack.ResultSize)
	if _, err := writer.Write(buf32); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}

	// BodySize
	binary.BigEndian.PutUint32(buf32, pack.BodySize)
	if _, err := writer.Write(buf32); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}

	// CertSize
	binary.BigEndian.PutUint16(buf16, pack.CertSize)
	if _, err := writer.Write(buf16); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}

	// Certificate
	if pack.CertSize > 0 {
		if _, err := writer.Write(pack.Certificate); err != nil {
			return errors.Errorf("write result file error: %s", err.Error())
		}
	}

	return nil
}

// PersistentJudgeResult 持久化评测记录
func PersistentJudgeResult(
	judgeResult *commonStructs.JudgeResult,
	options *persistence2.JudgeResultPersisOptions,
) error {
	fout, err := os.Create(options.OutFile)
	if err != nil {
		return errors.Errorf("create result file error: %s", err.Error())
	}
	defer fout.Close()

	if options.DigitalSign {
		if options.DigitalPEM.PublicKey == nil || options.DigitalPEM.PrivateKey == nil {
			return errors.Errorf("digital sign need public key and private key")
		}
	}

	resultBytes := utils.ObjectToJSONByte(judgeResult)

	bodyFile, err := mergeResultBinary(options, judgeResult)
	if err != nil {
		return err
	}

	bodyInfo, err := os.Stat(bodyFile)
	if err != nil {
		return err
	}

	certSize := 0 // 0 means disable cert
	var publicKeyRaw []byte
	if options.DigitalSign {
		certSize = len(options.DigitalPEM.PublicKeyRaw)
		publicKeyRaw = options.DigitalPEM.PublicKeyRaw
	}

	pack := JudgeResultPackage{
		Version:        1,
		Result:         resultBytes,
		ResultSize:     uint32(len(resultBytes)),
		BodySize:       uint32(bodyInfo.Size()),
		CertSize:       uint16(certSize),
		Certificate:    publicKeyRaw,
		CompressorType: options.CompressorType,
	}
	// Write Header
	err = writeFileHeaderAndResult(fout, pack)
	if err != nil {
		return err
	}

	// Write Signature
	fBody, err := os.Open(bodyFile)
	if err != nil {
		return err
	}

	hash, err := persistence2.SHA256Streams([]io.Reader{
		bytes.NewReader(resultBytes),
		fBody,
	})
	if err != nil {
		return err
	}
	_ = fBody.Close()
	if options.DigitalSign {
		hash, err = persistence2.RSA2048Sign(hash, options.DigitalPEM.PrivateKey)
		if err != nil {
			return err
		}
	}
	buf16 := make([]byte, 2)
	signSize := uint16(len(hash))
	// SignSize
	binary.BigEndian.PutUint16(buf16, signSize)
	if _, err := fout.Write(buf16); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}
	// Signature
	if _, err := fout.Write(hash); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}

	// Write Result and Body
	// 要注意先写入result，再写body，方便后续校验的时候直接顺序读取
	if _, err := fout.Write(pack.Result); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}
	fBody, err = os.Open(bodyFile)
	if err != nil {
		return err
	}
	defer fBody.Close()
	// Copy Body to fout
	if _, err := io.Copy(fout, fBody); err != nil {
		return errors.Errorf("write result file error: %s", err.Error())
	}
	// Clean
	_ = os.Remove(bodyFile)

	return nil
}
