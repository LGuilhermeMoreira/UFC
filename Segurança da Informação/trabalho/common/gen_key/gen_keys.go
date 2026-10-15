package genkey

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"golang.org/x/crypto/ssh"
)

// retorna a chave privada no formato pem e a chave publica no formato ssh
func GenKeys(filenameString string, createFileFlag bool) ([]byte, []byte, error) {
	// escolhe uma curva e gera uma chave privada a partir dessa curva
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, errors.New("falha ao gerar chave privada")
	}

	// faz o parse das chaves para bytes no formato DER
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, errors.New("falha ao converter a chave privada para bytes")
	}

	// cria um bloco PEM com os bytes no formato DER
	privateKeyPemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// faz o encode para pem
	var privateKeyPem bytes.Buffer
	if err := pem.Encode(&privateKeyPem, privateKeyPemBlock); err != nil {
		return nil, nil, errors.New("falha ao codificar a chave privada para PEM")
	}

	if createFileFlag {
		// cria o arquivo com filename e permissão 600
		if err := os.WriteFile(filenameString, privateKeyPem.Bytes(), 0600); err != nil {
			return nil, nil, errors.New("falha ao salvar a chave privada no ficheiro")
		}
	}

	// cria a chave SSH publica
	publicKeySSH, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, errors.New("falha ao criar a chave pública SSH")
	}
	// faz o marcha da chave publica para bytes
	publicKeySSHBytes := ssh.MarshalAuthorizedKey(publicKeySSH)

	publicKeyFileName := filenameString + ".pub"

	if createFileFlag {
		// cria o arquivo com filename + .pub,  com permissão 644
		if err := os.WriteFile(publicKeyFileName, publicKeySSHBytes, 0644); err != nil {
			return nil, nil, errors.New("falha ao salvar a chave pública no ficheiro")
		}
	}

	// retorna os bytes
	return privateKeyPem.Bytes(), publicKeySSHBytes, nil
}
