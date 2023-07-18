package persistence

import (
	"crypto/rsa"
)

// DigitalSignPEM 数字签名PEM证书内容
type DigitalSignPEM struct {
	PublicKey     *rsa.PublicKey
	PrivateKey    *rsa.PrivateKey
	PublicKeyRaw  []byte
	PrivateKeyRaw []byte
}

// CommonPersisOptions 公共的持久化选项
type CommonPersisOptions struct {
	DigitalSign bool            // 是否启用数字签名
	DigitalPEM  *DigitalSignPEM // PEM 数字证书
	OutFile     string          // 最终输出文件
}

// JudgeResultPersisOptions 评测记录的持久化选项
type JudgeResultPersisOptions struct {
	CommonPersisOptions
	CompressorType   uint8
	SessionDir       string
	SaveAcceptedData bool // 是否保存已经AC的数据
}

// ProblemPackageOptions 题目包的持久化选项
type ProblemPackageOptions struct {
	CommonPersisOptions
	ConfigFile string
	ConfigDir  string
}
