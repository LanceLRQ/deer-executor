package persistence

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/LanceLRQ/deer-executor/v3/executor/constants"
	"github.com/pkg/errors"
	"io"
	"os"
)

// sign body
func (pack *DeerPackageBase) signPackageBody(options *CommonPersisOptions) error {
	tBody, err := os.Open(options.TempFile)
	if err != nil {
		return err
	}
	defer tBody.Close()

	// Hash sign
	pack.Signature, err = SHA256Streams([]io.Reader{bufio.NewReader(tBody)})
	if err != nil {
		return err
	}
	return nil
}

// write body into package
func (pack *DeerPackageBase) writePackageBody(options *CommonPersisOptions, writer io.Writer) error {
	tBody, err := os.Open(options.TempFile)
	if err != nil {
		return err
	}
	defer tBody.Close()

	if _, err := io.Copy(writer, tBody); err != nil {
		return errors.Errorf("write problem file error: %s", err.Error())
	}
	return nil
}

// Do GPG signature
func (pack *DeerPackageBase) makeGPGSignature(options *CommonPersisOptions, writer io.Writer) error {
	buf16 := make([]byte, 2)

	pack.GPGCertSize = uint16(len(options.DigitalPEM.PublicKeyRaw))
	pack.GPGCertificate = options.DigitalPEM.PublicKeyRaw

	gpgHash, err := SHA256Streams([]io.Reader{
		bytes.NewReader(pack.PackageID.Bytes()),
		bytes.NewReader(pack.Signature),
	})
	if err != nil {
		return err
	}
	if pack.GPGCertSize <= 0 {
		return errors.Errorf("GPG certification is empty")
	}

	// CertSize
	binary.BigEndian.PutUint16(buf16, pack.GPGCertSize)
	if _, err := writer.Write(buf16); err != nil {
		return err
	}
	// Certificate
	if _, err := writer.Write(pack.GPGCertificate); err != nil {
		return err
	}
	// GPG sign
	gpgHash, err = RSA2048Sign(gpgHash, options.DigitalPEM.PrivateKey)
	if err != nil {
		return err
	}
	// SignSize
	signSize := uint16(len(gpgHash))
	binary.BigEndian.PutUint16(buf16, signSize)
	if _, err := writer.Write(buf16); err != nil {
		return err
	}
	// Signature
	if _, err := writer.Write(gpgHash); err != nil {
		return err
	}
	return nil
}

// Clean workspace temp file
func (pack *DeerPackageBase) cleanWorkspaceCommon(options *CommonPersisOptions) {
	if options.TempFile != "" {
		_ = os.Remove(options.TempFile)
	}
	options.TempFile = ""
}

// Write body chunk into file
func (pack *DeerPackageBase) writeBodyChunk(writer io.Writer, chunkType uint8, chunkSize uint64, chunkData io.Reader) error {
	buf16 := make([]byte, 2)
	buf32 := make([]byte, 4)
	buf64 := make([]byte, 8)

	errorFunc := func(key string, err error) error {
		return errors.Errorf("write problem chunk (%s) error: %s", key, err.Error())
	}

	flag := uint8(0)
	clen := uint8(1)        // 默认是用int16存储
	if chunkSize > 0xFFFF { // 大于2个字节，用int32
		clen = 2
	}
	if chunkSize > 0xFFFFFFFF { // 大于4个字节，用int64
		clen = 3
	}
	flag = (chunkType & MaskTypeNum) | ((clen << 5) & MaskTypeLen)
	if chunkSize > 0 {
		flag |= MaskTypeFlag
	}

	// Type Flag
	if _, err := writer.Write([]byte{flag}); err != nil {
		return errorFunc("type flag", err)
	}

	// Length
	var err error
	switch clen {
	case 1:
		binary.BigEndian.PutUint16(buf16, uint16(chunkSize))
		_, err = writer.Write(buf16)
	case 2:
		binary.BigEndian.PutUint32(buf32, uint32(chunkSize))
		_, err = writer.Write(buf32)
	case 3:
		binary.BigEndian.PutUint64(buf64, chunkSize)
		_, err = writer.Write(buf64)
	}
	if err != nil {
		return errorFunc("length", err)
	}

	// Content
	if chunkSize > 0 {
		if _, err := io.Copy(writer, chunkData); err != nil {
			return errorFunc("chunk content", err)
		}
	}

	return nil
}

// Create package header
func (pack *DeerPackageBase) createPackageHeader(writer io.Writer) error {
	buf16 := make([]byte, 2)

	errorFunc := func(key string, err error) error {
		return errors.Errorf("write problem header (%s) error: %s", key, err.Error())
	}

	// Magic Code
	binary.BigEndian.PutUint16(buf16, constants.ProblemPackageMagicCode)
	if _, err := writer.Write(buf16); err != nil {
		return errorFunc("magic", err)
	}
	// Version
	binary.BigEndian.PutUint16(buf16, pack.Version)
	if _, err := writer.Write(buf16); err != nil {
		return errorFunc("version", err)
	}

	// Commit Version
	binary.BigEndian.PutUint16(buf16, pack.CommitVersion)
	if _, err := writer.Write(buf16); err != nil {
		return errorFunc("commit version", err)
	}

	// Package ID
	if _, err := writer.Write(pack.PackageID.Bytes()); err != nil {
		return errorFunc("package id", err)
	}

	// Signature
	if _, err := writer.Write(pack.Signature); err != nil {
		return errorFunc("signature", err)
	}

	return nil
}
