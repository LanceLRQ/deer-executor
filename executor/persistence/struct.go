package persistence

import (
	"crypto/rsa"
	uuid "github.com/satori/go.uuid"
)

const (
	PackageChunkTypeConfig uint8 = iota + 1 // GPG Certificate
	PackageChunkTypeProject
	PackageChunkTypeJudgeResult
	PackageChunkTypeResultBody
)

const MaskTypeFlag uint8 = 0x80
const MaskTypeLen uint8 = 0x60
const MaskTypeNum uint8 = 0x1F

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
	TempFile    string          // 临时的文件，可清理
}

// JudgeResultPersisOptions 评测记录的持久化选项
type JudgeResultPersisOptions struct {
	CommonPersisOptions
	CompressorType   uint8
	SessionDir       string
	SaveAcceptedData bool // 是否保存已经AC的数据
}

type IDeerPackage interface {
	buildPackageBody(options *CommonPersisOptions) error
}

// DeerPackageBase deer executor data package binary struct definition
//
//	 ------------------------
//		header: |MAG|VER|CMT|PID|SIG|<GPCZ|GPGCert|GPSZ|GPGSign>|0x00 0x00|<body>
//		body: < | TYP | LEN | CNT | >
//		------------------------
//		Size
//		------------------------
//		header: | 2 | 2 | 2 | 16 | 32 | 2 | ... | 2 | ... |  (0x0)... |
//		body: | 1 | <2-4> | ... |
//		------------------------
//		TYPE
//		------------------------
//		|0|00|00000|
//		1 bit: Flag: continue or break
//		2-3 bit: <LEN> variable bytes length
//		4-8 bit: <CNT> type definition
type DeerPackageBase struct {
	IDeerPackage
	// --- header
	Version        uint16    // <VAR> Package Version
	CommitVersion  uint16    // <CMT> Commit Version
	PackageID      uuid.UUID // <PID> Package ID
	Signature      []byte    // <SIG> SHA256 Signature (sign body content)
	GPGCertSize    uint16    // (GPCZ) Public Certificate Size
	GPGCertificate []byte    // Public Certificate
	GPGSignSize    uint16    // (GPSZ) Signature Size
	GPGSignature   []byte    // Signature: SHA256(PackageId + Body)
	// --- body
	// ...
}
