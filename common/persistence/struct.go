package persistence

import (
	"crypto/rsa"
)

type DigitalSignPEM struct {
	PublicKey     *rsa.PublicKey
	PrivateKey    *rsa.PrivateKey
	PublicKeyRaw  []byte
	PrivateKeyRaw []byte
}

type CommonPersisOptions struct {
	DigitalSign bool            // 是否启用数字签名
	DigitalPEM  *DigitalSignPEM // PEM 数字证书
	OutFile     string          // 最终输出文件
}

type JudgeResultPersisOptions struct {
	CommonPersisOptions
	CompressorType   uint8
	SessionDir       string
	SaveAcceptedData bool // 是否保存已经AC的数据
}

type ProblemPackageOptions struct {
	CommonPersisOptions
	ConfigFile string
	ConfigDir  string
}
