package github

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

//! Github keys

func auxGetGithubKeys() ([]byte, error) {
	// faz a requisição para o github pegar as chaves
	resp, err := http.DefaultClient.Get("https://github.com/LGuilhermeMoreira.keys")
	if err != nil {
		return nil, err
	}
	// lê os bytes da resposta
	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// chega para ver se não veio nada
	if len(bts) == 0 {
		return nil, errors.New("nenhum dado encontrado.")
	}

	return bts, nil
}

// carrega a chave publica do cliente que está no github
func LoadClientECDSAPublicKey() (string, error) {
	bts, err := auxGetGithubKeys()
	if err != nil {
		return "", err
	}

	// converte os bytes para string
	allKeysString := string(bts)

	// separa as chaves
	keys := strings.Split(allKeysString, "\n")

	if len(keys) < 2 {
		return "", errors.New("mínimo necessário de chaves é 2\n1 para servidor e outra para cliente")
	}

	return keys[1], nil
}

// carrega a chave publica do servidor que está no github
func LoadServerECDSAPublicKey() (string, error) {
	bts, err := auxGetGithubKeys()
	if err != nil {
		return "", err
	}

	// converte os bytes para string
	allKeysString := string(bts)

	// separa as chaves
	keys := strings.Split(allKeysString, "\n")

	if len(keys) < 2 {
		return "", errors.New("mínimo necessário de chaves é 2\n1 para servidor e outra para cliente")
	}
	return keys[2], nil
}
