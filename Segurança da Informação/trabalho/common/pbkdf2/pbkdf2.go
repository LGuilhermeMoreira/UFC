package pbkdf2

import (
	"crypto/pbkdf2"
	"crypto/sha256"
)

const (
	HMACKeySize      = 32
	AESKeySize       = 32
	Pbkdf2Iterations = 100000
)

// baseado no sharedKey e no salto, gera a chave do AES e do HMAC
func DerivationFunc(sharedKey, salt []byte) ([]byte, []byte, error) {
	AESSalt := append(salt, []byte("aes")...)
	HMACSalt := append(salt, []byte("hmac")...)

	AESKey, err := pbkdf2.Key(sha256.New, string(sharedKey), AESSalt, Pbkdf2Iterations, AESKeySize)
	if err != nil {
		return nil, nil, err
	}
	HMACKey, err := pbkdf2.Key(sha256.New, string(sharedKey), HMACSalt, Pbkdf2Iterations, HMACKeySize)
	if err != nil {
		return nil, nil, err
	}

	return AESKey, HMACKey, nil
}

/*
pensei em implemnetar em gerar um chave só e divir para hmac e aes
, mas no slide estáva gerando 2 separadas.

*/
