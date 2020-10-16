package persistence

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path"
	"reflect"
)

func parseJudgeResultBinary(reader io.Reader ) (*JudgeResultPackage, error) {
	// 校验魔数
	magic := uint16(0)
	if err := binary.Read(reader, binary.BigEndian, &magic); err != nil {
		return nil, fmt.Errorf("read file error: %s", err.Error())
	}
	if magic != 0xB540 {
		return nil, fmt.Errorf("not deer-executor judge result file")
	}
	// 开始解析package
	pack := JudgeResultPackage{}
	if err := binary.Read(reader, binary.BigEndian, &pack.Version); err != nil {
		return nil, fmt.Errorf("read [version] error: %s", err.Error())
	}
	if err := binary.Read(reader, binary.BigEndian, &pack.CompressorType); err != nil {
		return nil, fmt.Errorf("read [compressor] type error: %s", err.Error())
	}
	if err := binary.Read(reader, binary.BigEndian, &pack.ResultSize); err != nil {
		return nil, fmt.Errorf("read [result size] error: %s", err.Error())
	}
	if err := binary.Read(reader, binary.BigEndian, &pack.BodySize); err != nil {
		return nil, fmt.Errorf("read [body size] error: %s", err.Error())
	}
	if err := binary.Read(reader, binary.BigEndian, &pack.CertSize); err != nil {
		return nil, fmt.Errorf("read [cert size] error: %s", err.Error())
	}
	// 如果有证书
	if pack.CertSize > 0 {
		pack.Certificate = make([]byte, pack.CertSize)
		if err := binary.Read(reader, binary.BigEndian, &pack.Certificate); err != nil {
			return nil, fmt.Errorf("read [cert public key] error: %s", err.Error())
		}
	}
	// 读取签名
	if err := binary.Read(reader, binary.BigEndian, &pack.SignSize); err != nil {
		return nil, fmt.Errorf("read [sign size] error: %s", err.Error())
	}
	pack.Signature = make([]byte, pack.SignSize)
	if err := binary.Read(reader, binary.BigEndian, &pack.Signature); err != nil {
		return nil, fmt.Errorf("read [signature] error: %s", err.Error())
	}
	// 读取Result
	pack.Result = make([]byte, pack.ResultSize)
	if err := binary.Read(reader, binary.BigEndian, &pack.Result); err != nil {
		return nil, fmt.Errorf("read [result] error: %s", err.Error())
	}
	// 理论上BodySize是多余的，剩下的都是body，这里就作为校验吧！
	tmpBodyFileName := uuid.NewV1().String() + ".tmp.gz"
	tmpBodyFilePath := path.Join("/tmp/", tmpBodyFileName)
	pack.BodyPackageFile = tmpBodyFilePath
	tmpBodyFile, err := os.Create(pack.BodyPackageFile)
	if err != nil {
		return nil, fmt.Errorf("create body package temp file error: %s", err.Error())
	}
	defer tmpBodyFile.Close()
	//var bodyReader io.Reader
	//// 如果使用了Gz
	//if pack.CompressorType == 1 {
	//	bodyReader, err = gzip.NewReader(tmpBodyFile)
	//	if err != nil {
	//		return nil, err
	//	}
	//} else {
	//	bodyReader = bufio.NewReader(tmpBodyFile)
	//}
	if _, err := io.Copy(tmpBodyFile, reader); err != nil {
		return nil, fmt.Errorf("write body package temp file error: %s", err.Error())
	}

	return &pack, nil
}

func validateJudgeResultPackage (pack *JudgeResultPackage) (bool, error) {
	// 打开临时文件
	tmpBodyFile, err := os.Open(pack.BodyPackageFile)
	if err != nil {
		return false, fmt.Errorf("open body package temp file error: %s", err.Error())
	}
	defer tmpBodyFile.Close()

	hash, err := SHA256Streams([]io.Reader{
		bytes.NewReader(pack.Result),
		tmpBodyFile,
	})
	if err != nil {
		return false, err
	}

	// 进行签名校验
	if pack.CertSize > 0 {
		publicKey, err := ReadAndParsePublicKey(pack.Certificate)
		if err != nil {
			return false, err
		}
		err = RSA2048Verify(hash, pack.Signature, publicKey)
		if err != nil {
			return false, err
		}
	} else {
		return reflect.DeepEqual(hash, pack.Signature), nil
	}
	return true, nil
}

func ReadJudgeResult(resultFile string) (*executor.JudgeResult, error) {
	rf, err := os.Open(resultFile)
	if err != nil {
		return nil, fmt.Errorf("open file (%s) error: %s", resultFile, err.Error())
	}
	reader := bufio.NewReader(rf)

	pack, err := parseJudgeResultBinary(reader)
	if err != nil {
		return nil, err
	}
	
	ok, err := validateJudgeResultPackage(pack)
	if !ok || err != nil {
		if err != nil {
			return nil, fmt.Errorf("validate package hash error: %s", err.Error())
		}
		return nil, fmt.Errorf("validate package hash error")
	}

	judgeResult := executor.JudgeResult{}
	executor.JSONBytesObject(pack.Result, &judgeResult)

	return &judgeResult, nil
}