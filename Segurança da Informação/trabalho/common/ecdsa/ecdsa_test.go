package ecdsa_test

import (
	"math/big"
	"testing"
	"trabalho-seguranca/common/ecdsa"
	"trabalho-seguranca/common/github"

	"github.com/stretchr/testify/assert"
)

func TestECDSA(t *testing.T) {
	chaveSecreta, err := ecdsa.GetPrivateKeyFromPem("../../keys/client_key")
	assert.NoError(t, err)
	assert.NotNil(t, chaveSecreta)

	chavePublicaString, err := github.LoadClientECDSAPublicKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, chavePublicaString)

	chavePublica, err := ecdsa.GetPublicKeyFromSSH([]byte(chavePublicaString))
	assert.NoError(t, err)
	assert.NotNil(t, chavePublica)

	sign, err := ecdsa.Sing(chaveSecreta, big.NewInt(2), "client_user")
	assert.NoError(t, err)
	assert.NotEmpty(t, sign)

	ok := ecdsa.Verify(chavePublica, sign, big.NewInt(2), "client_user")
	assert.True(t, ok)
}
