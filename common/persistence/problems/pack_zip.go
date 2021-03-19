package problems

import (
	"archive/zip"
	"github.com/LanceLRQ/deer-executor/v2/common/persistence"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func packZipFile(rootPath string, targetPath string) error {
	fout, err := os.Create(targetPath)
	if err != nil {
		return errors.Errorf("create problem package file error: %s", err.Error())
	}
	defer fout.Close()

	zipWriter := zip.NewWriter(fout)
	defer zipWriter.Close()

	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, rootPath+"/")

		// 排除根目录
		if header.Name == rootPath {
			return nil
		}
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		// 排除bin目录
		if strings.HasPrefix(header.Name, "bin/") {
			return nil
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

// 执行题目数据表打包操作（打包成zip版本)
func PackProblemsAsZip(options *persistence.ProblemPackageOptions) error {
	// 这边没法支持所有内容的校验了，只能给config签名。
	if options.DigitalSign {
		if options.DigitalPEM.PublicKey == nil || options.DigitalPEM.PrivateKey == nil {
			return errors.Errorf("digital sign need public key and private key")
		}
	}

	pubkeyFileName := path.Join(options.ConfigDir, ".gpg")
	signFileName := path.Join(options.ConfigDir, ".sign")

	// clean meta file
	if _, err := os.Stat(pubkeyFileName); os.IsExist(err) {
		_ = os.Remove(pubkeyFileName)
	}
	if _, err := os.Stat(signFileName); os.IsExist(err) {
		_ = os.Remove(signFileName)
	}

	defer func() {
		_ = os.Remove(pubkeyFileName)
		_ = os.Remove(signFileName)
	}()

	// Cretea Public Key
	if options.DigitalSign {
		gpgFile, err := os.Create(pubkeyFileName)
		if err != nil {
			return err
		}
		if _, err = gpgFile.Write(options.DigitalPEM.PublicKeyRaw); err != nil {
			return err
		}
		_ = gpgFile.Close()
	}

	// Create Signature (only for configFile)
	fBody, err := os.Open(options.ConfigFile)
	if err != nil {
		return err
	}

	hash, err := persistence.SHA256Streams([]io.Reader{fBody})

	// GPG signature
	if options.DigitalSign {
		hash, err = persistence.RSA2048Sign(hash, options.DigitalPEM.PrivateKey)
		if err != nil {
			return err
		}
	}

	// Write to .sign file
	signFile, err := os.Create(signFileName)
	if err != nil {
		return err
	}
	if _, err = signFile.Write(hash); err != nil {
		return err
	}
	_ = signFile.Close()

	// Package Zip
	err = packZipFile(options.ConfigDir, options.OutFile)

	return err
}
