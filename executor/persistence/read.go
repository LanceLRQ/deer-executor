package persistence

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"github.com/LanceLRQ/deer-executor/v3/executor/constants"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
	"reflect"
)

// parse package's header
func (pack *DeerPackageBase) parseDeerPackageHeader(reader io.Reader) error {
	errorFunc := func(key string, err error) error {
		return errors.Errorf("read package file (%s) error: %s", key, err.Error())
	}

	// check magic number
	magic := uint32(0)
	if err := binary.Read(reader, binary.BigEndian, &magic); err != nil {
		return errorFunc("magic number", err)
	}
	if magic != constants.DeerPackageMagicCode {
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

// validate package body
func (pack *DeerPackageBase) validatePackage(curHash []byte) error {
	isOk := reflect.DeepEqual(curHash, pack.Signature)
	if !isOk {
		return errors.Errorf("validate package signature error")
	}
	if pack.GPGCertSize > 0 {
		gpgHash, err := SHA256Streams([]io.Reader{
			bytes.NewReader(pack.PackageID.Bytes()),
			bytes.NewReader(pack.Signature),
		})
		if err != nil {
			return err
		}
		elist, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(pack.GPGCertificate))
		if err != nil {
			return err
		}
		if len(elist) < 1 {
			return errors.Errorf("GPG key error")
		}
		publicKey := elist[0].PrimaryKey.PublicKey.(*rsa.PublicKey)
		err = RSA2048Verify(gpgHash, pack.GPGSignature, publicKey)
		if err != nil {
			return errors.Errorf("validate package error: %s", err.Error())
		}
	}
	return nil
}

// parse package common func
func (pack *DeerPackageBase) parseDeerPackage(doValidate bool) error {
	// Open pack file
	fp, err := os.Open(pack.presistFilePath)
	if err != nil {
		return errors.Errorf("open file (%s) error: %s", pack.presistFilePath, err.Error())
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)

	// parse header
	err = pack.parseDeerPackageHeader(reader)
	if err != nil {
		return err
	}
	if !doValidate {
		return nil
	}

	// Hash check
	bodyHash, bodySize, err := SHA256StreamsWithCount([]io.Reader{
		reader,
	})
	if err != nil {
		return nil
	}

	pack.presistHeaderSize = pack.presistFileSize - bodySize

	// Validate
	err = pack.validatePackage(bodyHash)
	if err != nil {
		return nil
	}

	return nil
}

func (pack *DeerPackageBase) walkDeerPackageBody(items []interface{}, callback func(uint8, int64, *io.SectionReader) (int64, error)) error {
	// Open pack file
	fp, err := os.Open(pack.presistFilePath)
	if err != nil {
		return errors.Errorf("open file (%s) error: %s", pack.presistFilePath, err.Error())
	}
	defer fp.Close()

	// Seek to body
	_, err = fp.Seek(pack.presistHeaderSize, io.SeekStart)
	if err != nil {
		return err
	}

	for {
		typeId, cLength, err := pack.readPackageBodyChunkHeader(fp)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if items != nil && typeId > 0 {
			if _, ok := isInArray(typeId, items); !ok {
				// next <cLength> bytes
				_, err = fp.Seek(cLength, io.SeekCurrent)
				if err != nil {
					return err
				}
				continue
			}
		}
		// get current file pointer position
		fpos, err := fp.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		// read operation must in callbak
		n, err := callback(typeId, cLength, io.NewSectionReader(fp, fpos, cLength))
		if err != nil {
			return err
		}
		if n < cLength { // if not read all bytes, skip
			_, err := fp.Seek(cLength-n, io.SeekCurrent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (pack *DeerPackageBase) readPackageBodyChunkHeader(file *os.File) (uint8, int64, error) {
	flags := make([]byte, 1)
	if _, err := file.Read(flags); err != nil {
		return 0, 0, err
	}
	flag := flags[0]
	isContinue := (flag&MaskTypeFlag)>>7 == 1
	chunkSizeLen := (flag & MaskTypeLen) >> 5
	chunkSize := uint64(0)
	chunkType := flag & MaskTypeNum
	if !isContinue {
		return 0, 0, nil
	}
	switch chunkSizeLen {
	case 1:
		l := uint16(0)
		if err := binary.Read(file, binary.BigEndian, &l); err != nil {
			return 0, 0, err
		}
		chunkSize = uint64(l)
	case 2:
		l := uint32(0)
		if err := binary.Read(file, binary.BigEndian, &l); err != nil {
			return 0, 0, err
		}
		chunkSize = uint64(l)
	case 3:
		if err := binary.Read(file, binary.BigEndian, &chunkSize); err != nil {
			return 0, 0, err
		}
	}
	return chunkType, int64(chunkSize), nil
}
