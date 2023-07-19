package problems

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/LanceLRQ/deer-executor/v3/executor/constants"
	"github.com/LanceLRQ/deer-executor/v3/executor/persistence"
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

func NewProblemPackage(conf *commonStructs.JudgeConfiguration) ProblemPackage {
	return ProblemPackage{
		Version:        2,
		PackageID:      uuid.NewV4(),
		ProblemConfigs: conf,
	}
}

// WritePackageFile 打包题目数据到文件
func (pack *ProblemPackage) WritePackageFile(options *persistence.ProblemPackageOptions) error {
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
	err = pack.signPackageBody(options)
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

	// GPG
	if options.DigitalSign {
		err = pack.makeGPGSignature(options, foutWriter)
		if err != nil {
			return errors.Errorf("sign problem package file error: %s", err.Error())
		}
	}

	// Header end with [0 0]
	if _, err := foutWriter.Write([]byte{0x0, 0x0}); err != nil {
		return err
	}

	// Write body
	err = pack.writePackageBody(options, foutWriter)
	if err != nil {
		return err
	}

	return nil
}

// 合并题目项目文件
func (pack *ProblemPackage) mergeFilesBinary(options *persistence.ProblemPackageOptions) error {
	zipFile, err := os.Create(options.BodyFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(options.ConfigDir, func(zpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过配置文件
		if zpath == options.ConfigFile {
			return nil
		}
		// 跳过配置文件夹
		if zpath == options.ConfigDir {
			return nil
		}
		// 跳过二进制文件夹
		binDir := path.Join(options.ConfigDir, "bin")
		if strings.HasPrefix(zpath, binDir) {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(zpath, options.ConfigDir+"/")
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

// 对body进行签名
func (pack *ProblemPackage) signPackageBody(options *persistence.ProblemPackageOptions) error {
	tBody, err := os.Open(options.TempFile)
	if err != nil {
		return err
	}
	defer tBody.Close()

	// Hash sign
	pack.Signature, err = persistence.SHA256Streams([]io.Reader{bufio.NewReader(tBody)})
	if err != nil {
		return err
	}
	return nil
}

// 把body写到数据包
func (pack *ProblemPackage) writePackageBody(options *persistence.ProblemPackageOptions, writer io.Writer) error {
	tBody, err := os.Open(options.TempFile)
	if err != nil {
		return err
	}
	defer tBody.Close()

	if _, err := io.Copy(writer, tBody); err != nil {
		return errors.Errorf("write problem file error: %s", err.Error())
	}
	return nil
}

// 清理数据
func (pack *ProblemPackage) cleanWorkspace(options *persistence.ProblemPackageOptions) {
	if options.TempFile != "" {
		_ = os.Remove(options.TempFile)
	}
	if options.BodyFile != "" {
		_ = os.Remove(options.BodyFile)
	}
	options.BodyFile = ""
	options.TempFile = ""
}

// GPG signature
func (pack *ProblemPackage) makeGPGSignature(options *persistence.ProblemPackageOptions, writer io.Writer) error {
	buf16 := make([]byte, 2)

	pack.GPGCertSize = uint16(len(options.DigitalPEM.PublicKeyRaw))
	pack.GPGCertificate = options.DigitalPEM.PublicKeyRaw

	gpgHash, err := persistence.SHA256Streams([]io.Reader{
		bytes.NewReader(pack.PackageID.Bytes()),
		bytes.NewReader(pack.Signature),
	})
	if err != nil {
		return err
	}
	if pack.GPGCertSize <= 0 {
		return errors.Errorf("GPG certification is empty")
	}

	// CertSize
	binary.BigEndian.PutUint16(buf16, pack.GPGCertSize)
	if _, err := writer.Write(buf16); err != nil {
		return err
	}
	// Certificate
	if _, err := writer.Write(pack.GPGCertificate); err != nil {
		return err
	}
	// GPG sign
	gpgHash, err = persistence.RSA2048Sign(gpgHash, options.DigitalPEM.PrivateKey)
	if err != nil {
		return err
	}
	// SignSize
	signSize := uint16(len(gpgHash))
	binary.BigEndian.PutUint16(buf16, signSize)
	if _, err := writer.Write(buf16); err != nil {
		return err
	}
	// Signature
	if _, err := writer.Write(gpgHash); err != nil {
		return err
	}
	return nil
}

// 构建body
func (pack *ProblemPackage) buildPackageBody(options *persistence.ProblemPackageOptions) error {
	options.TempFile = fmt.Sprintf("%s.tmp", options.OutFile)
	tempOut, err := os.Create(options.TempFile)
	if err != nil {
		return errors.Errorf("create problem package body file error: %s", err.Error())
	}
	defer tempOut.Close()
	tempOutWriter := bufio.NewWriter(tempOut)

	// config
	configBytes := utils.ObjectToJSONByte(pack.ProblemConfigs)
	err = pack.writeBodyChunk(tempOutWriter, ProblemPackageBodyTypeConfig, uint64(len(configBytes)), bytes.NewBuffer(configBytes))
	if err != nil {
		return err
	}

	// content file
	options.BodyFile = fmt.Sprintf("%s.zip", options.OutFile)
	err = pack.mergeFilesBinary(options)
	if err != nil {
		return err
	}

	fBody, err := os.Open(options.BodyFile)
	if err != nil {
		return err
	}
	bodyInfo, err := fBody.Stat()
	if err != nil {
		return err
	}
	err = pack.writeBodyChunk(tempOutWriter, ProblemPackageBodyTypeContent, uint64(bodyInfo.Size()), fBody)
	if err != nil {
		return err
	}

	return nil
}

// 生成一个body chunk
func (pack *ProblemPackage) writeBodyChunk(writer io.Writer, chunkType uint8, chunkSize uint64, chunkData io.Reader) error {
	buf16 := make([]byte, 2)
	buf32 := make([]byte, 4)
	buf64 := make([]byte, 8)

	errorFunc := func(key string, err error) error {
		return errors.Errorf("write problem chunk (%s) error: %s", key, err.Error())
	}

	flag := uint8(0)
	clen := uint8(1)        // 默认是用int16存储
	if chunkSize > 0xFFFF { // 大于2个字节，用int32
		clen = 2
	}
	if chunkSize > 0xFFFFFFFF { // 大于4个字节，用int64
		clen = 3
	}
	flag = (chunkType & MaskTypeNum) | ((clen << 5) & MaskTypeLen)
	if chunkSize > 0 {
		flag |= MaskTypeFlag
	}

	// Type Flag
	if _, err := writer.Write([]byte{flag}); err != nil {
		return errorFunc("type flag", err)
	}

	// Length
	var err error
	switch clen {
	case 1:
		binary.BigEndian.PutUint16(buf16, uint16(chunkSize))
		_, err = writer.Write(buf16)
	case 2:
		binary.BigEndian.PutUint32(buf32, uint32(chunkSize))
		_, err = writer.Write(buf32)
	case 3:
		binary.BigEndian.PutUint64(buf64, chunkSize)
		_, err = writer.Write(buf64)
	}
	if err != nil {
		return errorFunc("length", err)
	}

	// Content
	if chunkSize > 0 {
		if _, err := io.Copy(writer, chunkData); err != nil {
			return errorFunc("chunk content", err)
		}
	}

	return nil
}

// 构建一个header
func (pack *ProblemPackage) createPackageHeaderBuffer(writer io.Writer) error {
	buf16 := make([]byte, 2)

	errorFunc := func(key string, err error) error {
		return errors.Errorf("write problem header (%s) error: %s", key, err.Error())
	}

	// Magic Code
	binary.BigEndian.PutUint16(buf16, constants.ProblemPackageMagicCode)
	if _, err := writer.Write(buf16); err != nil {
		return errorFunc("magic", err)
	}
	// Version
	binary.BigEndian.PutUint16(buf16, pack.Version)
	if _, err := writer.Write(buf16); err != nil {
		return errorFunc("version", err)
	}

	// Commit Version
	binary.BigEndian.PutUint16(buf16, pack.CommitVersion)
	if _, err := writer.Write(buf16); err != nil {
		return errorFunc("commit version", err)
	}

	// Package ID
	if _, err := writer.Write(pack.PackageID.Bytes()); err != nil {
		return errorFunc("package id", err)
	}

	// Signature
	if _, err := writer.Write(pack.Signature); err != nil {
		return errorFunc("signature", err)
	}

	return nil
}
