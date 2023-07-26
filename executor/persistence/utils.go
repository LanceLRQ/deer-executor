package persistence

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// SHA256String 对文本执行SHA256运算
func SHA256String(body string) ([]byte, error) {
	return SHA256Bytes([]byte(body))
}

// SHA256Bytes 对byte数组执行SHA256运算
func SHA256Bytes(body []byte) ([]byte, error) {
	buf := bytes.NewBuffer(body)
	return SHA256Streams([]io.Reader{buf})
}

// SHA256Streams 对文件流执行SHA256运算
func SHA256Streams(streams []io.Reader) ([]byte, error) {
	hash := sha256.New()
	for _, stream := range streams {
		if _, err := io.Copy(hash, stream); err != nil {
			return nil, err
		}
	}
	ret := hash.Sum(nil)
	//hex.EncodeToString(ret)
	return ret, nil
}

// SHA256StreamsWithCount sha256 and return size
func SHA256StreamsWithCount(streams []io.Reader) ([]byte, int64, error) {
	counter := int64(0)
	hash := sha256.New()
	for _, stream := range streams {
		if c, err := io.Copy(hash, stream); err != nil {
			return nil, 0, err
		} else {
			counter += c
		}
	}
	ret := hash.Sum(nil)
	//hex.EncodeToString(ret)
	return ret, counter, nil
}

// RSA2048SignString 对文本执行RSA2048运算
func RSA2048SignString(body string, privateKey *rsa.PrivateKey) ([]byte, error) {
	return RSA2048SignBytes([]byte(body), privateKey)
}

// RSA2048SignBytes 对byte数组执行RSA2048运算
func RSA2048SignBytes(body []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash, err := SHA256Bytes(body)
	if err != nil {
		return nil, err
	}
	return RSA2048Sign(hash, privateKey)
}

// RSA2048Sign 用RSA2048运算签名
func RSA2048Sign(hash []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

// RSA2048VerifyString 用RSA2048运算校验文本签名
func RSA2048VerifyString(body string, sign []byte, publicKey *rsa.PublicKey) error {
	return RSA2048VerifyBytes([]byte(body), sign, publicKey)
}

// RSA2048VerifyBytes 用RSA2048运算校验byte数组签名
func RSA2048VerifyBytes(body []byte, sign []byte, publicKey *rsa.PublicKey) error {
	hash, err := SHA256Bytes(body)
	if err != nil {
		return err
	}
	return RSA2048Verify(hash, sign, publicKey)
}

// RSA2048Verify 用RSA2048运算校验签名
func RSA2048Verify(hash []byte, sign []byte, publicKey *rsa.PublicKey) error {
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, sign)
}

// ReadAndParsePublicKey 解析PEM证书公钥
func ReadAndParsePublicKey(cert []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(cert)
	if block == nil {
		return nil, errors.New("read public key error")
	}
	//publicKey, err := x509.ParseCertificate(block.Bytes)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return publicKey.(*rsa.PublicKey), nil
	//return publicKey.PublicKey.(*rsa.PublicKey), nil
}

// ReadAndParsePrivateKey 解析PEM证书私钥
func ReadAndParsePrivateKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("read private key error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// ReadPemFile 读取PEM证书
func ReadPemFile(filename string) ([]byte, error) {
	pemBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return pemBytes, nil
}

// GetDigitalPEMFromFile 从文件读取PEM证书
func GetDigitalPEMFromFile(publicKeyFile string, privateKeyFile string) (*DigitalSignPEM, error) {
	publicKey, err := ReadPemFile(publicKeyFile)
	if err != nil {
		return nil, err
	}
	privateKey, err := ReadPemFile(privateKeyFile)
	if err != nil {
		return nil, err
	}
	dp := GetDigitalPEM(publicKey, privateKey)
	return dp, nil
}

// GetDigitalPEM 获取PEM数据
func GetDigitalPEM(publicKey []byte, privateKey []byte) *DigitalSignPEM {
	dp := DigitalSignPEM{}
	if publicKey != nil {
		dp.PublicKeyRaw = publicKey
		pk, err := ReadAndParsePublicKey(publicKey)
		dp.PublicKey = pk
		if err != nil {
			dp.PublicKey = nil
		}
	}
	if privateKey != nil {
		dp.PrivateKeyRaw = privateKey
		pk, err := ReadAndParsePrivateKey(privateKey)
		dp.PrivateKey = pk
		if err != nil {
			dp.PrivateKey = nil
		}
	}
	return &dp
}

// Gets 实现gets方法
func Gets(reader io.Reader) string {
	buf := make([]byte, 16, 16)
	rd := bufio.NewReader(reader)
	for {
		t, err := rd.ReadByte()
		if err != nil {
			return string(buf)
		}
		buf = append(buf, t)
		if t == 13 || t == 0 {
			break
		}
	}
	return string(buf)
}

// GetPublicKeyArmorBytes 读取GPG证书信息
func GetPublicKeyArmorBytes(entity *openpgp.Entity) ([]byte, error) {
	pathname := path.Join("/tmp/" + uuid.NewV4().String() + ".key")
	fp, err := os.Create(pathname)
	if err != nil {
		return nil, err
	}
	defer os.Remove(pathname)

	w, err := armor.Encode(fp, openpgp.PublicKeyType, nil)
	if err != nil {
		return nil, err
	}
	err = entity.Serialize(w)
	if err != nil {
		return nil, err
	}
	_ = w.Close()
	_ = fp.Close()

	data, err := ioutil.ReadFile(pathname)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetArmorPublicKey 从文件读取GPG证书公钥信息
func GetArmorPublicKey(gpgKeyFile string, passphrase []byte) (*DigitalSignPEM, error) {
	if gpgKeyFile == "" {
		return nil, errors.Errorf("please set GPG key file path")
	}
	// Open key
	keyRingReader, err := os.Open(gpgKeyFile)
	if err != nil {
		return nil, err
	}
	// Read GPG Keys
	elist, err := openpgp.ReadArmoredKeyRing(keyRingReader)
	if err != nil {
		return nil, err
	}
	if len(elist) < 1 {
		return nil, errors.Errorf("file has no GPG key")
	}
	gpgKey := elist[0].PrivateKey
	if gpgKey.Encrypted {
		if len(passphrase) == 0 {
			passphrase, err = gopass.GetPasswdPrompt("please input passphrase of key> ", true, os.Stdin, os.Stdout)
			if err != nil {
				return nil, err
			}
		}
		err = gpgKey.Decrypt(passphrase)
		if err != nil {
			return nil, err
		}
	}
	publicKeyArmor, err := GetPublicKeyArmorBytes(elist[0])
	if err != nil {
		return nil, err
	}
	return &DigitalSignPEM{
		PrivateKey:   gpgKey.PrivateKey.(*rsa.PrivateKey),
		PublicKeyRaw: publicKeyArmor,
		PublicKey:    elist[0].PrimaryKey.PublicKey.(*rsa.PublicKey),
	}, nil
}

// FileNotFoundError file not found error
type FileNotFoundError struct {
	FileName string
}

// Error get error string
func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("file (%s) not found", e.FileName)
}

// IsFileNotFoundError check is file not found error
func IsFileNotFoundError(e error) bool {
	_, ok := e.(FileNotFoundError)
	return ok
}

// NewFileNotFoundError new a file not found error
func NewFileNotFoundError(fileName string) error {
	return FileNotFoundError{FileName: fileName}
}

// UnZip do unzip
func UnZip(zipArchive *zip.ReadCloser, destDir string) error {
	for _, f := range zipArchive.File {
		err := WalkZip(f, destDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnZipReader do unzip (for *zip.Reader)
func UnZipReader(zipArchive *zip.Reader, destDir string) error {
	for _, f := range zipArchive.File {
		err := WalkZip(f, destDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// WalkZip walk in zip file
func WalkZip(f *zip.File, destDir string) error {
	fpath := filepath.Join(destDir, f.Name)
	if f.FileInfo().IsDir() {
		_ = os.MkdirAll(fpath, os.ModePerm)
	} else {
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
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
	return nil
}

// FindInZip find file in zip exactly.
func FindInZip(zipArchive *zip.ReadCloser, fileName string) (*io.ReadCloser, *zip.File, error) {
	var fileResult io.ReadCloser
	finded := false
	var fileInfo *zip.File
	for _, file := range zipArchive.File {
		fileInfo = file
		if fileName == file.Name {
			finded = true
			var err error
			fileResult, err = file.Open()
			if err != nil {
				return nil, nil, err
			}
		}
	}
	if !finded {
		return nil, nil, NewFileNotFoundError(fileName)
	}
	return &fileResult, fileInfo, nil
}

func isInArray(target interface{}, arr []interface{}) (int, bool) {
	for i, value := range arr {
		if target == value {
			return i, true
		}
	}
	return -1, false
}
