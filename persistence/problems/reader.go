package problems

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/LanceLRQ/deer-executor/persistence"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"syscall"
)

// 解析判题结果
func parseProblemPackageBinary(reader io.Reader, unpackBody bool) (*ProblemPackage, error) {
	// 校验魔数
	magic := uint16(0)
	if err := binary.Read(reader, binary.BigEndian, &magic); err != nil {
		return nil, fmt.Errorf("read file error: %s", err.Error())
	}
	if magic != persistence.ProblemPackageMagicCode {
		return nil, fmt.Errorf("not deer-executor problem package file")
	}
	// 开始解析package
	pack := ProblemPackage{}
	if err := binary.Read(reader, binary.BigEndian, &pack.Version); err != nil {
		return nil, fmt.Errorf("read [version] error: %s", err.Error())
	}
	if err := binary.Read(reader, binary.BigEndian, &pack.CommitVersion); err != nil {
		return nil, fmt.Errorf("read [version] error: %s", err.Error())
	}
	if err := binary.Read(reader, binary.BigEndian, &pack.ConfigSize); err != nil {
		return nil, fmt.Errorf("read [config size] error: %s", err.Error())
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
	// 读取Config
	pack.Configs = make([]byte, pack.ConfigSize)
	if err := binary.Read(reader, binary.BigEndian, &pack.Configs); err != nil {
		return nil, fmt.Errorf("read [config] error: %s", err.Error())
	}
	if unpackBody {
		// 理论上BodySize是多余的，剩下的都是body，这里就作为校验吧！
		tmpBodyFileName := uuid.NewV1().String() + ".tmp.zip"
		tmpBodyFilePath := path.Join("/tmp/", tmpBodyFileName)
		pack.BodyPackageFile = tmpBodyFilePath
		tmpBodyFile, err := os.Create(pack.BodyPackageFile)
		if err != nil {
			return nil, fmt.Errorf("create body package temp file error: %s", err.Error())
		}
		defer tmpBodyFile.Close()
		if _, err := io.Copy(tmpBodyFile, reader); err != nil {
			return nil, fmt.Errorf("write body package temp file error: %s", err.Error())
		}
	}

	return &pack, nil
}

// 校验判题结果数据包
func validateProblemPackage (pack *ProblemPackage) (bool, error) {
	// 打开临时文件
	tmpBodyFile, err := os.Open(pack.BodyPackageFile)
	if err != nil {
		return false, fmt.Errorf("open body package temp file error: %s", err.Error())
	}
	defer tmpBodyFile.Close()

	hash, err := persistence.SHA256Streams([]io.Reader{
		bytes.NewReader(pack.Configs),
		tmpBodyFile,
	})
	if err != nil {
		return false, err
	}
	// 进行签名校验
	if pack.CertSize > 0 {
		// Read GPG Keys
		elist, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(pack.Certificate))
		if err != nil {
			return false, err
		}
		if len(elist) < 1 {
			return false, fmt.Errorf("GPG key error")
		}
		publicKey := elist[0].PrimaryKey.PublicKey.(*rsa.PublicKey)
		err = persistence.RSA2048Verify(hash, pack.Signature, publicKey)
		if err != nil {
			return false, err
		}
	} else {
		return reflect.DeepEqual(hash, pack.Signature), nil
	}
	return true, nil
}

func readProblemPackage(problemFile string, validate bool, unpackBody bool) (*ProblemPackage, error){
	fp, err := os.Open(problemFile)
	if err != nil {
		return nil, fmt.Errorf("open file (%s) error: %s", problemFile, err.Error())
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)
	pack, err := parseProblemPackageBinary(reader, unpackBody)
	if err != nil {
		return nil, err
	}
	if validate && unpackBody {
		ok, err := validateProblemPackage(pack)
		if !ok || err != nil {
			if err != nil {
				return nil, fmt.Errorf("validate package hash error: %s", err.Error())
			}
			return nil, fmt.Errorf("validate package hash error")
		}
	}
	return pack, nil
}

func unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}


// 读取题目信息
func ReadProblemInfo(problemFile string, unpackBody bool, workDir string) (*executor.JudgeSession, error) {

	pack, err := readProblemPackage(problemFile, true, unpackBody)
	if err != nil {
		return nil, err
	}
	session := executor.JudgeSession{}
	executor.JSONBytesObject(pack.Configs, &session.JudgeConfig)

	if unpackBody {
		err = unzip(pack.BodyPackageFile, workDir)
		if err != nil {
			return nil, err
		}
		configFile := path.Join(workDir, "problem.json")
		fp, err := os.Create(configFile)
		if err != nil {
			return nil, err
		}
		defer fp.Close()
		_, err = fp.Write(pack.Configs)
		if err != nil {
			return nil, err
		}
		session.ConfigFile = configFile
	}

	return &session, nil
}

// 读取题目携带的GPG信息
func ReadProblemGPGInfo(problemFile string) (string, error) {
	pack, err := readProblemPackage(problemFile, false, false)
	if err != nil {
		return "", err
	}
	if pack.CertSize == 0 {
		return "no GPG public key", nil
	} else {
		elist, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(pack.Certificate))
		if err != nil {
			return "", err
		}
		if len(elist) < 1 {
			return "", fmt.Errorf("GPG key error")
		}
		rel := ""
		for _, identify := range elist[0].Identities {
			rel += identify.Name + "\n"
		}
		return rel, nil
	}
}

// 判断是否是题目包
func IsProblemPackage (filePath string) (bool, error) {
	fp, err := os.OpenFile(filePath, os.O_RDONLY | syscall.O_NONBLOCK, 0)
	if err != nil {
		return false, fmt.Errorf("open file error")
	}
	defer fp.Close()

	var magic uint16 = 0
	err = binary.Read(fp, binary.BigEndian, &magic)
	if err != nil {
		return false, err
	}

	return  magic == persistence.ProblemPackageMagicCode, nil
}