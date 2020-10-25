package judge_result

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/LanceLRQ/deer-executor/persistence"
	"github.com/LanceLRQ/deer-executor/persistence/judge_result"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path"
)



func writeFileHeaderAndResult (writer io.Writer, pack ProblemPackage) error {
	buf16 := make([]byte, 2)
	buf32 := make([]byte, 4)

	// magic
	binary.BigEndian.PutUint16(buf16, persistence.JudgeResultMagicCode)
	if _, err := writer.Write(buf16); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}

	// Version
	binary.BigEndian.PutUint16(buf16, pack.Version)
	if _, err := writer.Write(buf16); err != nil {
		return fmt.Errorf("write result file error: %s", err.Error())
	}

	// ResultSize
	binary.BigEndian.PutUint32(buf32, pack.ConfigSize)
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

func PackProblems(
	session *executor.JudgeSession,
	options ProblemPersisOptions,
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

	configBytes := executor.ObjectToJSONByte(session)

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

	pack := ProblemPackage{
		Version: 1,
		Configs: configBytes,
		ConfigSize: uint32(len(configBytes)),
		BodySize: uint32(bodyInfo.Size()),
		CertSize: uint16(certSize),
		Certificate: options.DigitalPEM.PublicKeyRaw,
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
		bytes.NewReader(configBytes),
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
	if _, err := fout.Write(pack.Configs); err != nil {
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