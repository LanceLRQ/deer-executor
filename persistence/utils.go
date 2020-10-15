package persistence

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
)

/* SHA256 */

func SHA256String(body string) ([]byte, error) {
	return SHA256Bytes([]byte(body))
}

func SHA256Bytes(body []byte) ([]byte, error) {
	buf := bytes.NewBuffer(body)
	return SHA256Stream(buf)
}

func SHA256Stream(body io.Reader) ([]byte, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, body); err != nil {
		return nil, err
	}
	ret := hash.Sum(nil)
	//hex.EncodeToString(ret)
	return ret, nil
}

/* RSA Sign */

func RSA2048SignString(body string, key []byte) ([]byte, error) {
	return RSA2048SignBytes([]byte(body), key)
}

func RSA2048SignBytes(body []byte, key []byte) ([]byte, error) {
	buf := bytes.NewBuffer(body)
	return RSA2048SignStream(buf, key)
}

func RSA2048SignStream(body io.Reader, key []byte) ([]byte, error) {
	hash, err := SHA256Stream(body)
	if err != nil { return nil, err }
	// read private key
	privateKey, err := ReadAndParsePrivateKey(key)
	if err != nil { return nil, err }
	// sign with private key
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

/* RSA Verify */

func RSA2048VerifyString(body string, sign []byte, cert []byte) error {
	return RSA2048VerifyBytes([]byte(body), sign, cert)
}

func RSA2048VerifyBytes(body []byte, sign []byte, cert []byte) error {
	buf := bytes.NewBuffer(body)
	return RSA2048VerifyStream(buf, sign, cert)
}

func RSA2048VerifyStream(body io.Reader, sign []byte, cert []byte) error {
	hash, err := SHA256Stream(body)
	if err != nil { return err }
	// read public key
	publicKey, err := ReadAndParsePublicKey(cert)
	if err != nil { return err }
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, sign)
}

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

func ReadPemFile(filename string) ([]byte, error) {
	pemBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return pemBytes, nil
}

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