package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	// flag para gerar nome do arquivo
	fileName := flag.String("fn", "id_ecdsa", "Nome base para gerar os arquivos de chave SSH.")
	flag.Parse()

	// gera a chave privada e utiliza a curva P256
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gerar chave privada: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Chaves geradas na memória com sucesso.")

	// converte a chave privada para bytes.
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao converter chave privada: %v\n", err)
		os.Exit(1)
	}

	privateKeyPem := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// salva o arquivo da chave privada
	privateKeyFile, err := os.Create(*fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo de chave privada: %v\n", err)
		os.Exit(1)
	}
	defer privateKeyFile.Close()

	// define as permissões do arquivo da chave privada
	if err := os.Chmod(*fileName, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao definir permissões do arquivo: %v\n", err)
		os.Exit(1)
	}

	if err := pem.Encode(privateKeyFile, privateKeyPem); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao salvar chave privada no arquivo: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Chave privada SSH salva em: %s\n", *fileName)

	// pega a chave publica
	publicKey := privateKey.Public()

	// converte a chave para formato ssh
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar chave pública SSH: %v\n", err)
		os.Exit(1)
	}

	// faz o marshal da chave
	authorizedKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)

	// salva o arquivo
	publicKeyFileName := *fileName + ".pub"
	publicKeyFile, err := os.Create(publicKeyFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo de chave pública: %v\n", err)
		os.Exit(1)
	}
	defer publicKeyFile.Close()

	_, err = publicKeyFile.Write(authorizedKeyBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao salvar chave pública no arquivo: %v\n", err)
		os.Exit(1)
	}
}
