package problems

import (
	commonStructs "github.com/LanceLRQ/deer-executor/v3/executor/structs"
	uuid "github.com/satori/go.uuid"
)

/* ProblemPackage 题目包-数据结构
------------------------
header: |MAG|VER|CMT|PID|SIG|<GPCZ|GPGCert|GPSZ|GPGSign>|0x00 0x00|<body>
body: < | TYP | LEN | CNT | >
------------------------
Size
------------------------
header: | 2 | 2 | 2 | 16 | 32 | 2 | ... | 2 | ... |  (0x0)... |
body: | 1 | <2-4> | ... |
------------------------
TYPE
------------------------
|0|00|00000|
1 bit: Flag: continue or break
2-3 bit: <LEN> variable bytes length
4-8 bit: <CNT> type definition
*/

const MaskTypeFlag uint8 = 0x80
const MaskTypeLen uint8 = 0x60
const MaskTypeNum uint8 = 0x1F

const (
	ProblemPackageBodyTypeConfig uint8 = iota + 1 // GPG Certificate
	ProblemPackageBodyTypeContent
)

type ProblemPackage struct {
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
	gpgCertificate      []byte // GPG Public Certificate (optional) [type: 0x1]
	gpgSignature        []byte // GPG Public Certificate (optional) [type: 0x2]
	problemConfigsBytes []byte // Problem Configs JSON [type: 0x3]
	ProblemBodyTempFile string // Problem package temp file [type: 0x4]

	// META
	ProblemConfigs *commonStructs.JudgeConfiguration
}
