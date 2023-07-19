package problems

//
//import (
//	"bytes"
//	"encoding/binary"
//	persistence "github.com/LanceLRQ/deer-executor/v3/executor/persistence"
//	commonStructs "github.com/LanceLRQ/deer-executor/v3/executor/structs"
//	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
//	"github.com/pkg/errors"
//	"io"
//	"os"
//)
//
//func writeFileHeader(writer io.Writer, pack ProblemPackage) error {
//	buf16 := make([]byte, 2)
//	buf32 := make([]byte, 4)
//
//	// ConfigSize
//	binary.BigEndian.PutUint32(buf32, pack.ConfigSize)
//	if _, err := writer.Write(buf32); err != nil {
//		return errors.Errorf("write problem file error: %s", err.Error())
//	}
//
//	// BodySize
//	binary.BigEndian.PutUint32(buf32, pack.BodySize)
//	if _, err := writer.Write(buf32); err != nil {
//		return errors.Errorf("write problem file error: %s", err.Error())
//	}
//
//	// CertSize
//	binary.BigEndian.PutUint16(buf16, pack.CertSize)
//	if _, err := writer.Write(buf16); err != nil {
//		return errors.Errorf("write problem file error: %s", err.Error())
//	}
//
//	// Certificate
//	if pack.CertSize > 0 {
//		if _, err := writer.Write(pack.Certificate); err != nil {
//			return errors.Errorf("write problem file error: %s", err.Error())
//		}
//	}
//
//	return nil
//}
//
//// PackProblems 执行题目数据表打包操作
//func PackProblems(
//	configuration *commonStructs.JudgeConfiguration,
//	options *persistence.ProblemPackageOptions,
//) error {
//
//	if options.DigitalSign {
//		if options.DigitalPEM.PublicKey == nil || options.DigitalPEM.PrivateKey == nil {
//			return errors.Errorf("digital sign need public key and private key")
//		}
//	}
//
//	fout, err := os.Create(options.OutFile)
//	if err != nil {
//		return errors.Errorf("create problem package file error: %s", err.Error())
//	}
//	defer fout.Close()
//
//	configBytes := utils.ObjectToJSONByte(configuration)
//
//	bodyFile, err := mergeFilesBinary(options)
//	if err != nil {
//		return err
//	}
//
//	bodyInfo, err := os.Stat(bodyFile)
//	if err != nil {
//		return err
//	}
//
//	certSize := 0 // 0 means disable cert
//	var publicKeyRaw []byte
//	if options.DigitalSign {
//		certSize = len(options.DigitalPEM.PublicKeyRaw)
//		publicKeyRaw = options.DigitalPEM.PublicKeyRaw
//	}
//
//	pack := ProblemPackage{
//		Version:       1,
//		CommitVersion: 1,
//		Configs:       configBytes,
//		ConfigSize:    uint32(len(configBytes)),
//		BodySize:      uint32(bodyInfo.Size()),
//		CertSize:      uint16(certSize),
//		Certificate:   publicKeyRaw,
//	}
//	// Write Header
//	err = writeFileHeader(fout, pack)
//	if err != nil {
//		return err
//	}
//
//	// Write Signature
//	fBody, err := os.Open(bodyFile)
//	if err != nil {
//		return err
//	}
//
//	hash, err := persistence.SHA256Streams([]io.Reader{
//		bytes.NewReader(configBytes),
//		fBody,
//	})
//	if err != nil {
//		return err
//	}
//	_ = fBody.Close()
//	// GPG signature
//	if options.DigitalSign {
//		hash, err = persistence.RSA2048Sign(hash, options.DigitalPEM.PrivateKey)
//		if err != nil {
//			return err
//		}
//	}
//	buf16 := make([]byte, 2)
//	signSize := uint16(len(hash))
//	// SignSize
//	binary.BigEndian.PutUint16(buf16, signSize)
//	if _, err := fout.Write(buf16); err != nil {
//		return errors.Errorf("write problem file error: %s", err.Error())
//	}
//	// Signature
//	if _, err := fout.Write(hash); err != nil {
//		return errors.Errorf("write problem file error: %s", err.Error())
//	}
//
//	// Write configs and Body
//	if _, err := fout.Write(pack.ProblemConfigs); err != nil {
//		return errors.Errorf("write problem file error: %s", err.Error())
//	}
//	fBody, err = os.Open(bodyFile)
//	if err != nil {
//		return err
//	}
//	defer fBody.Close()
//	// Copy Body to fout
//	if _, err := io.Copy(fout, fBody); err != nil {
//		return errors.Errorf("write problem file error: %s", err.Error())
//	}
//
//	// Clean
//	_ = os.Remove(bodyFile)
//
//	return nil
//
//}
