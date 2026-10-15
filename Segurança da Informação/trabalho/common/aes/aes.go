package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
)

// essa função utiliza AES-CBC, recebe como parametro a chave e mensagem em texto claro.
// É adicionando padding por causa do CBC, durante a função é gerado o IV que é retornado
func Encrypt(aesKey, plaintext []byte) (ciphertext, iv []byte, err error) {
	// cria o bloco
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, err
	}

	// adiciona o padding no texto claro
	paddedText := addPKCS7Padding(plaintext, block.BlockSize())

	// gera o iv
	iv = make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}

	// cria o buf do texto cifrado
	ciphertext = make([]byte, len(paddedText))

	// organiza para cifrar a mensagem com o bloco e o iv
	cbcEncrypter := cipher.NewCBCEncrypter(block, iv)

	// carrega o buffer com texto claro com padding criptografado
	cbcEncrypter.CryptBlocks(ciphertext, paddedText)

	return ciphertext, iv, nil
}

// descriptografa uma mensagem usando AES-CBC.
func Decrypt(aesKey, iv, ciphertext []byte) ([]byte, error) {
	// cria o bloco
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// verifica o tamanho do bloco para ver se veio correto
	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, errors.New("ciphertext não é múltiplo do tamanho do bloco")
	}

	// cria o buffer para texto decriptogrado com padding
	decrypted := make([]byte, len(ciphertext))

	// prepara para decifrar
	cbcDecrypter := cipher.NewCBCDecrypter(block, iv)

	// carrega o buffer com o texto decifrado
	cbcDecrypter.CryptBlocks(decrypted, ciphertext)

	// remove o padding
	plaintext, err := removePKCS7Padding(decrypted)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

/*
adiciona o padding no texto claro
ex:
OLA MUNDO tem 9 bytes, mas o tamanho do bloco é 16 bytes
então precisamos completar os 9 para ter 16.
texto claro com padding.
[O, L, Á,  , M, U, N, D, O, 7, 7, 7, 7, 7, 7, 7].
*/
func addPKCS7Padding(data []byte, blockSize int) []byte {
	// faz o calculo do padding
	padding := blockSize - (len(data) % blockSize)
	// cria um array com o padding.
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	// adiciona o padding no texto clarro
	return append(data, padText...)
}

// remove o padding do texto claro
func removePKCS7Padding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("dados de padding vazios")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return nil, fmt.Errorf("padding inválido: %d", padding)
	}
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, errors.New("bytes de padding inconsistentes")
		}
	}
	return data[:len(data)-padding], nil
}
