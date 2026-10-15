package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net"
	"os"
	"trabalho-seguranca/common/aes"
	"trabalho-seguranca/common/df"
	"trabalho-seguranca/common/ecdsa"
	genkey "trabalho-seguranca/common/gen_key"
	"trabalho-seguranca/common/github"
	"trabalho-seguranca/common/hmac"
	"trabalho-seguranca/common/network"
	"trabalho-seguranca/common/path"
	"trabalho-seguranca/common/pbkdf2"

	"github.com/joho/godotenv"
)

var enableLogs bool
var salt string

func init() {
	flag.BoolVar(&enableLogs, "log", false, "Habilitar logs detalhados")
	flag.Parse()
	godotenv.Load()

	if !enableLogs {
		log.SetOutput(io.Discard)
	}
	salt = os.Getenv("SALT")
}

func main() {
	fmt.Println("[SERVIDOR] Iniciando servidor...")

	content, err := os.ReadFile(path.ServerKeyPath)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao ler a chave privada do servidor: %v", err)
	}

	if len(content) == 0 {
		log.Println("[SERVIDOR] Nenhuma chave encontrada. Gerando nova chave privada ECDSA...")
		_, _, errGenKey := genkey.GenKeys(path.ServerKeyPath, true)
		if errGenKey != nil {
			log.Fatalf("[ERRO] Falha ao gerar chave: %v", errGenKey)
		}
	}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("[ERRO] Falha ao iniciar servidor: %v", err)
	}
	log.Println("[SERVIDOR] Aguardando conexões na porta 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[ERRO] Falha ao aceitar conexão: %v", err)
			continue
		}
		log.Println("[SERVIDOR] Nova conexão aceita.")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("[SERVIDOR] Nova conexão aceita. Iniciando handshake...")

	clientPublicKeyString, err := github.LoadClientECDSAPublicKey()
	if err != nil {
		log.Printf("[ERRO] Falha ao carregar chave pública do cliente: %v", err)
		return
	}
	log.Println("[SERVIDOR] Chave pública do cliente ECDSA carregada com sucesso.")

	clientPublicKey, err := ecdsa.GetPublicKeyFromSSH([]byte(clientPublicKeyString))
	if err != nil {
		log.Printf("[ERRO] Falha ao converter chave pública do cliente: %v", err)
		return
	}

	ABytes, err := network.ReadField(conn)
	if err != nil {
		log.Printf("[ERRO] Falha ao ler chave DH do cliente: %v", err)
		return
	}
	A := big.NewInt(0).SetBytes(ABytes)
	log.Printf("[SERVIDOR] Chave Diffie-Hellman do cliente recebida: %v", A)

	clientSign, err := network.ReadField(conn)
	if err != nil {
		log.Printf("[ERRO] Falha ao ler assinatura do cliente: %v", err)
		return
	}

	clientName, err := network.ReadField(conn)
	if err != nil {
		log.Printf("[ERRO] Falha ao ler nome do cliente: %v", err)
		return
	}
	log.Printf("[SERVIDOR] Nome do cliente: %s", string(clientName))

	ok := ecdsa.Verify(clientPublicKey, clientSign, A, string(clientName))
	if !ok {
		log.Println("[ERRO] Falha na verificação da assinatura do cliente.")
		return
	}
	log.Println("[SERVIDOR] Assinatura do cliente verificada com sucesso.")

	serverName := fmt.Sprintf("server_%v", rand.Int())

	b, err := df.GenerateSecretKey(df.Prime)
	if err != nil {
		log.Printf("[ERRO] Falha ao gerar chave secreta DF do servidor: %v", err)
		return
	}
	log.Printf("[SERVIDOR] Chave secreta DF gerada: %v", b)

	B, err := df.GenerateSharedKeys(df.Prime, df.Generator, b)
	if err != nil {
		log.Printf("[ERRO] Falha ao gerar chave pública DF: %v", err)
		return
	}
	log.Printf("[SERVIDOR] Chave pública DF gerada: %v", B)

	privateKey, err := ecdsa.GetPrivateKeyFromPem(path.ServerKeyPath)
	if err != nil {
		log.Printf("[ERRO] Falha ao obter chave privada do servidor: %v", err)
		return
	}

	sign, err := ecdsa.Sing(privateKey, B, serverName)
	if err != nil {
		log.Printf("[ERRO] Falha ao assinar chave DF: %v", err)
		return
	}
	log.Printf("[SERVIDOR] Assinatura da chave DF gerada.")

	data, err := network.BuildHandshakeMessage(B.Bytes(), sign, serverName)
	if err != nil {
		log.Printf("[ERRO] Falha ao construir mensagem de handshake: %v", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Printf("[ERRO] Falha ao enviar dados para o cliente: %v", err)
		return
	}
	log.Println("[SERVIDOR] Handshake enviado com sucesso.")

	S, err := df.GenerateSharedSecretKey(df.Prime, A, b)
	if err != nil {
		log.Printf("[ERRO] Falha ao gerar chave secreta compartilhada: %v", err)
		return
	}
	log.Printf("[SERVIDOR] Chave Secreta Compartilhada Diffie-Hellman: %v", S)

	AESKey, HMACKey, err := pbkdf2.DerivationFunc(S.Bytes(), []byte(salt))

	if err != nil {
		log.Printf("[ERRO] Falha ao derivar as chaves: %v", err)
		return
	}

	log.Println("[SERVIDOR] Derivação concluída. Aguardando mensagens do cliente...")

	for {
		HMACClientTag, err := network.ReadField(conn)
		if err != nil {
			if err == io.EOF {
				log.Printf("[SERVIDOR] Cliente %s desconectado.", string(clientName))
			} else {
				log.Printf("[ERRO] Falha ao ler tag HMAC enviado pelo cliente: %v", err)
			}
			break
		}
		log.Println("[SERVIDOR] Tag HMAC recebida.")

		iv, err := network.ReadField(conn)
		if err != nil {
			log.Printf("[ERRO] Falha ao ler IV enviado pelo cliente: %v", err)
			break
		}
		log.Println("[SERVIDOR] IV recebido.")

		cipherText, err := network.ReadField(conn)
		if err != nil {
			log.Printf("[ERRO] Falha ao ler cipher text enviado pelo cliente: %v", err)
			break
		}
		log.Printf("[SERVIDOR] Cipher text recebido.")

		okTagHMAC := hmac.Verify(HMACKey, HMACClientTag, iv, cipherText)
		if !okTagHMAC {
			log.Println("[ERRO] HMAC é inválido. Ignorando mensagem.")
			continue
		}
		log.Println("[SERVIDOR] HMAC validado.")

		text, err := aes.Decrypt(AESKey, iv, cipherText)
		if err != nil {
			log.Printf("[ERRO] Erro ao decifrar o texto cifrado: %v", err)
			continue
		}

		fmt.Printf("%s enviou: %s -> %s\n", string(clientName), string(cipherText), string(text))
	}
}
