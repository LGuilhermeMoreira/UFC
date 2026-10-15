package network

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

/*
eu escrevo um array contendo o tamanho e o dado para todos os valores passados por parametro
ex:
BuildHandshakeMessage(mensagem,ola) -> [2,men,sagem,1,ola]
*/
func BuildHandshakeMessage(sharedKey, sign []byte, username string) ([]byte, error) {
	// crio um buffer
	var buf bytes.Buffer

	// escreve o tamanho da chave compartilhada
	err := binary.Write(&buf, binary.BigEndian, uint32(len(sharedKey)))
	if err != nil {
		return nil, err
	}
	// escrevo a mensagem compartilhada
	_, err = buf.Write(sharedKey)
	if err != nil {
		return nil, err
	}
	// escreve o tamanho da assinatura
	err = binary.Write(&buf, binary.BigEndian, uint32(len(sign)))
	if err != nil {
		return nil, err
	}
	// escreve a assinatura
	_, err = buf.Write(sign)
	if err != nil {
		return nil, err
	}
	// faz o parse para bytes
	usernameBytes := []byte(username)
	// escreve o tamanho do username
	err = binary.Write(&buf, binary.BigEndian, uint32(len(usernameBytes)))
	if err != nil {
		return nil, err
	}
	// escreve o username
	buf.Write(usernameBytes)

	return buf.Bytes(), nil
}

/*
eu escrevo um array contendo o tamanho e o dado para todos os valores passados por parametro
ex:
BuildHandshakeMessage(mensagem,ola) -> [2,men,sagem,1,ola]
*/
func BuildMassage(HMACTag, IVAES, cipherText []byte) ([]byte, error) {
	var buf bytes.Buffer
	// escreto o tamanho da tag hmac
	err := binary.Write(&buf, binary.BigEndian, uint32(len(HMACTag)))
	if err != nil {
		return nil, err
	}
	// escrevo a tag hmac
	_, err = buf.Write(HMACTag)
	if err != nil {
		return nil, err
	}

	// escrevo o tamanho do iv aes
	err = binary.Write(&buf, binary.BigEndian, uint32(len(IVAES)))
	if err != nil {
		return nil, err
	}
	// escreve o iv aes
	_, err = buf.Write(IVAES)
	if err != nil {
		return nil, err
	}

	// escreve o tamanho de cipherText
	err = binary.Write(&buf, binary.BigEndian, uint32(len(cipherText)))
	if err != nil {
		return nil, err
	}

	// escreve o cipherText
	_, err = buf.Write(cipherText)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ReadField(conn net.Conn) ([]byte, error) {
	lenBuf := make([]byte, 4)
	// lê o tamanho do proximo campo
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf)

	field := make([]byte, length)
	// le o campo
	if _, err := io.ReadFull(conn, field); err != nil {
		return nil, err
	}

	// retorna o campo em bytes
	return field, nil
}
