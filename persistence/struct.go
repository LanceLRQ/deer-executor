package persistence

import "crypto/rsa"

type DigitalSignPEM struct {
	PublicKey 			*rsa.PublicKey
	PrivateKey 			*rsa.PrivateKey
	PublicKeyRaw 		[]byte
	PrivateKeyRaw 		[]byte
}
