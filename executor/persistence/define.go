package persistence

import (
	"crypto/rsa"
	uuid "github.com/satori/go.uuid"

	commonStructs "github.com/LanceLRQ/deer-executor/v3/executor/structs"
)

const (
	PackageTypeProblem uint16 = iota + 1
	PackageTypeJudgeResult
)

var PackageTypeNameMap = map[uint16]string{
	PackageTypeProblem:     "Problem Project",
	PackageTypeJudgeResult: "Judge Result",
}

const (
	PackageChunkTypeConfig uint8 = iota + 1
	PackageChunkTypeProject
	PackageChunkTypeJudgeResult
	PackageChunkTypeJudgeResultData
)

const MaskTypeFlag uint8 = 0x80
const MaskTypeLen uint8 = 0x60
const MaskTypeNum uint8 = 0x1F

// DigitalSignPEM gpg cert data
type DigitalSignPEM struct {
	PublicKey     *rsa.PublicKey
	PrivateKey    *rsa.PrivateKey
	PublicKeyRaw  []byte
	PrivateKeyRaw []byte
}

// CommonPersisOptions common persis options
type CommonPersisOptions struct {
	DigitalSign bool            // enable digital sign?
	DigitalPEM  *DigitalSignPEM // PEM certification data
	OutFile     string          // package result output file
	TempFile    string          // temporary file
}

// ProblemProjectPersisOptions problem persis options
type ProblemProjectPersisOptions struct {
	CommonPersisOptions
	ConfigFile      string // for skip the config file when zipping
	ProjectDir      string // problem project work dir
	ProblemBodyFile string
}

// JudgeResultPersisOptions judge resule persis options
type JudgeResultPersisOptions struct {
	CommonPersisOptions
	SessionDir          string
	JudgeResultDataFile string
	SaveAcceptedData    bool
}

type IDeerPackage interface {
	getCommonPersisOptions() (*CommonPersisOptions, error)
	buildPackageBody() error
	cleanWorkspace()
}

// DeerPackageBase deer executor data package binary struct definition
//
//	 ------------------------
//		header: |MAG|VER|CMV|PKT|PID|SIG|<GPCZ|GPGCert|GPSZ|GPGSign>|0x00 0x00|<body>
//		body: < | TYP | LEN | CNT | >
//		------------------------
//		Size
//		------------------------
//		header: | 4 | 2 | 2 | 16 | 32 | 2 | ... | 2 | ... |  (0x0)... |
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
	CommitVersion  uint16    // <CMV> Commit Version
	PackageType    uint16    // <PKT> Package Type
	PackageID      uuid.UUID // <PID> Package ID
	Signature      []byte    // <SIG> SHA256 Signature (sign body content)
	GPGCertSize    uint16    // (GPCZ) Public Certificate Size
	GPGCertificate []byte    // Public Certificate
	GPGSignSize    uint16    // (GPSZ) Signature Size
	GPGSignature   []byte    // Signature: SHA256(PackageId + Body)
	// --- body
	// ...
	// -- meta
	packageSize       int64
	presistOptions    interface{}
	presistHeaderSize int64
	presistFilePath   string
	presistFileSize   int64
}

// ProblemProjectPackage problem project package entity
type ProblemProjectPackage struct {
	DeerPackageBase
	// --- internal
	ProblemConfigs *commonStructs.JudgeConfiguration
}

type JudgeResultPackage struct {
	DeerPackageBase
	// -- body
	JudgeResult *commonStructs.JudgeResult
}
