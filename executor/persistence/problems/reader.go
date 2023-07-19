package problems

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/LanceLRQ/deer-executor/v3/executor/constants"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path"
)

// UnPackPackageFile 解析题目包
func UnPackPackageFile(packFile string, workDir string, headerOnly bool, doValidate bool) (*ProblemPackage, error) {
	// Open pack file
	fp, err := os.Open(packFile)
	if err != nil {
		return nil, errors.Errorf("open file (%s) error: %s", packFile, err.Error())
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)

	pack := ProblemPackage{}

	// parse header
	err = pack.parseProblemPackageHeader(reader)
	if err != nil {
		return nil, err
	}

	if !doValidate && headerOnly {
		return &pack, nil
	}

	tmpFileName := fmt.Sprintf("%s.tmp", uuid.NewV1().String(), ".tmp")
	pack.BodyPackageFile = path.Join(workDir, tmpFileName)

	// validation

	return &pack, nil
}

// Release packages

// Parse header
func (pack *ProblemPackage) parseProblemPackageHeader(reader io.Reader) error {

	errorFunc := func(key string, err error) error {
		return errors.Errorf("read package file (%s) error: %s", key, err.Error())
	}

	// check magic number
	magic := uint16(0)
	if err := binary.Read(reader, binary.BigEndian, &magic); err != nil {
		return errorFunc("magic number", err)
	}
	if magic != constants.ProblemPackageMagicCode {
		return errors.Errorf("not deer-executor problem package file")
	}
	if err := binary.Read(reader, binary.BigEndian, &pack.Version); err != nil {
		return errorFunc("version", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &pack.CommitVersion); err != nil {
		return errorFunc("commit version", err)
	}
	// uuid
	uuidBytes := make([]byte, 16)
	if err := binary.Read(reader, binary.BigEndian, &uuidBytes); err != nil {
		return errorFunc("package id", err)
	}
	pkid, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return err
	}
	pack.PackageID = pkid
	// hash
	hashBytes := make([]byte, 32)
	if err := binary.Read(reader, binary.BigEndian, &hashBytes); err != nil {
		return errorFunc("signature", err)
	}
	pack.Signature = hashBytes
	// Parse next 2 bytes
	dividerBytes := uint16(0)
	exit := false
	parsedSign := false
	for !exit {
		if err := binary.Read(reader, binary.BigEndian, &dividerBytes); err != nil {
			return err
		}
		if dividerBytes > 0 && !parsedSign { // not divider
			pack.GPGCertSize = dividerBytes
			pack.GPGCertificate = make([]byte, pack.GPGCertSize)
			if err := binary.Read(reader, binary.BigEndian, &pack.GPGCertificate); err != nil {
				return errorFunc("gpg public key", err)
			}
			// read sign
			if err := binary.Read(reader, binary.BigEndian, &pack.GPGSignSize); err != nil {
				return errorFunc("gpg sign size", err)
			}
			pack.GPGSignature = make([]byte, pack.GPGSignSize)
			if err := binary.Read(reader, binary.BigEndian, &pack.GPGSignature); err != nil {
				return errorFunc("gpg sign", err)
			}
			parsedSign = true
		} else {
			exit = true
		}
	}

	return nil
}

// validate package
func (pack *ProblemPackage) validatePackage(reader io.Reader) error {

	return nil
}
