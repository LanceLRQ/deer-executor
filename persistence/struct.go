package persistence

import (
	"crypto/rsa"
)

type DigitalSignPEM struct {
	PublicKey 			*rsa.PublicKey
	PrivateKey 			*rsa.PrivateKey
	PublicKeyRaw 		[]byte
	PrivateKeyRaw 		[]byte
}


type CommonPersisOptions struct {
	DigitalSign    bool
	DigitalPEM     *DigitalSignPEM
	OutFile        string
}

type JudgeResultPersisOptions struct {
	CommonPersisOptions
	CompressorType uint8
}
