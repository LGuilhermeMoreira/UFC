package ecdsa

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"

	"golang.org/x/crypto/ssh"
)

//! ecdsa

// pega a chave privada local utilizando o path dela como parametro
func GetPrivateKeyFromPem(path string) (*ecdsa.PrivateKey, error) {
	// pega os bytes do arquivo de chave privada
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// faz o decode dos bytes em bloco no formato PEM
	block, _ := pem.Decode(keyBytes)

	// se der algum bucho durante o decode, o block vem nil
	if block == nil {
		return nil, errors.New("erro ao fazer decode do bloco PEM da chave privada")
	}

	// olha se é do tipo correto
	if block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("bloco PEM não é do tipo 'EC PRIVATE KEY'")
	}

	// faz o parse para chave privada
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// pega a chave publica utilizando os bytes da chave como parametro
func GetPublicKeyFromSSH(publicKeyBts []byte) (*ecdsa.PublicKey, error) {
	// faz o parse para o tipo sshKey
	sshKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyBts)
	if err != nil {
		return nil, err
	}

	// faz o parse para o tipo ssh.CryptoPublicKey
	cryptoPubKey, ok := sshKey.(ssh.CryptoPublicKey)
	if !ok {
		return nil, errors.New("a chave parseada não é uma chave criptográfica")
	}

	// faz o parse de ssh.CryptoPublicKey para crypto.PublicKey e dps para *ecdsa.PublicKey
	ecdsaKey, ok := cryptoPubKey.CryptoPublicKey().(*ecdsa.PublicKey)

	if !ok {
		return nil, errors.New("a chave pública não é do tipo ECDSA")
	}

	return ecdsaKey, nil
}

func Sing(privateKey *ecdsa.PrivateKey, sharedKey *big.Int, userName string) ([]byte, error) {
	// transformo tudo em bytes
	dataToSign := sharedKey.Bytes()
	dataToSign = append(dataToSign, []byte(userName)...)

	// gera o hash
	hash := sha256.Sum256(dataToSign)

	// assina com a chave privada
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])

	if err != nil {
		return nil, err
	}

	return signature, nil
}

func Verify(publicKey *ecdsa.PublicKey, signature []byte, sharedKey *big.Int, userName string) bool {
	// transforma tudo em bytes
	dataToVerify := sharedKey.Bytes()
	dataToVerify = append(dataToVerify, []byte(userName)...)

	// gera o hash
	hash := sha256.Sum256(dataToVerify)

	// verifica
	return ecdsa.VerifyASN1(publicKey, hash[:], signature)
}
