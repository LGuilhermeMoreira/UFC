package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net"
	"os"
	"strings"
	"trabalho-seguranca/common/aes"
	"trabalho-seguranca/common/df"
	"trabalho-seguranca/common/ecdsa"
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
	godotenv.Load()
	flag.BoolVar(&enableLogs, "log", false, "Habilitar logs detalhados")
	flag.Parse()
	if !enableLogs {
		log.SetOutput(io.Discard)
	}
	salt = os.Getenv("SALT")
}

func main() {
	num := rand.Int()
	clientName := fmt.Sprintf("client_%v", num)
	log.Printf("[CLIENTE] Nome gerado: %s", clientName)

	a, err := df.GenerateSecretKey(df.Prime)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao gerar chave secreta DF: %v", err)
	}
	log.Printf("[CLIENTE] Chave secreta DF gerada: %v", a)

	A, err := df.GenerateSharedKeys(df.Prime, df.Generator, a)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao gerar chave compartilhada DF: %v", err)
	}
	log.Printf("[CLIENTE] Chave pública DF: %v", A)

	clientPrivateKey, err := ecdsa.GetPrivateKeyFromPem(path.ClientKeyPath)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao carregar chave privada do cliente: %v", err)
	}

	sign, err := ecdsa.Sing(clientPrivateKey, A, clientName)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao assinar dados: %v", err)
	}
	log.Println("[CLIENTE] Assinatura gerada com sucesso.")

	bts, err := network.BuildHandshakeMessage(A.Bytes(), sign, clientName)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao montar mensagem de handshake: %v", err)
	}

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("[ERRO] Falha ao conectar no servidor: %v", err)
	}
	defer conn.Close()
	log.Println("[CLIENTE] Conectado ao servidor.")

	_, err = conn.Write(bts)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao enviar dados para o servidor: %v", err)
	}
	log.Println("[CLIENTE] Handshake enviado.")

	serverPublicKeyString, err := github.LoadServerECDSAPublicKey()
	if err != nil {
		log.Fatalf("[ERRO] Falha ao carregar chave pública do servidor: %v", err)
	}
	log.Println("[CLIENTE] Chave pública do servidor carregada.")

	serverPublicKey, err := ecdsa.GetPublicKeyFromSSH([]byte(serverPublicKeyString))
	if err != nil {
		log.Fatalf("[ERRO] Falha ao converter chave pública do servidor: %v", err)
	}

	BBytes, err := network.ReadField(conn)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao ler chave DF do servidor: %v", err)
	}
	B := big.NewInt(0).SetBytes(BBytes)
	log.Printf("[CLIENTE] Chave DF recebida do servidor: %v", B)

	serverSign, err := network.ReadField(conn)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao ler assinatura do servidor: %v", err)
	}

	serverName, err := network.ReadField(conn)
	if err != nil {
		log.Fatalf("[ERRO] Falha ao ler nome do servidor: %v", err)
	}
	log.Printf("[CLIENTE] Nome do servidor: %s", string(serverName))

	ok := ecdsa.Verify(serverPublicKey, serverSign, B, string(serverName))
	if !ok {
		log.Fatalf("[ERRO] Assinatura do servidor inválida!")
	}
	log.Println("[CLIENTE] Assinatura do servidor verificada com sucesso.")

	S, err := df.GenerateSharedSecretKey(df.Prime, B, a)

	if err != nil {
		log.Fatalf("[ERRO] Falha ao gerar chave secreta compartilhada: %v", err)
	}

	log.Printf("[CLIENTE] Chave Secreta Compartilhada Diffie-Hellman: %v", S)

	log.Println("[CLIENTE] Derivando as chaves AES e HMAC a partir da chave secreta compartilhada.")

	AESKey, HMACKey, err := pbkdf2.DerivationFunc(S.Bytes(), []byte(salt))

	if err != nil {
		log.Fatalf("[ERRO] Falha ao derivar as chaves: %v", err)
	}

	log.Println("[CLIENTE] Derivação concluida.")

	for {
		fmt.Println("Caso queria parar de enviar mensagem para o serivodor, digite 'exit'")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("input: ")
		message, _ := reader.ReadString('\n')

		// Remove o caractere de nova linha no final
		message = strings.TrimSpace(message)
		log.Printf("[CLIENTE] Mensagem a ser enviada para o servidor: %v", string(message))

		if message == "exit" {
			break
		}

		cipherText, iv, err := aes.Encrypt(AESKey, []byte(message))
		if err != nil {
			log.Fatalf("[ERRO] Falha ao cifrar a mensagem: %v", err)
		}

		HMACTag := hmac.Calculate(HMACKey, iv, cipherText)

		log.Printf("[CLIENTE] Texto cifrado e HMAC calculado..")

		bts, err = network.BuildMassage(HMACTag, iv, cipherText)
		if err != nil {
			log.Fatalf("[ERRO] Falha gerar dados da mensagem: %v", err)
		}

		log.Printf("[CLIENTE] Enviado texto cifrado, HMAC e IV para o servidor")

		conn.Write(bts)
	}

}
