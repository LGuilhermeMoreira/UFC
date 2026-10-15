package pbkdf2_test

import (
	"testing"

	"trabalho-seguranca/common/pbkdf2"

	"github.com/stretchr/testify/assert"
)

func TestDerivationFunc(t *testing.T) {
	sharedKey := []byte("chave-compartilhada")
	salt := []byte("sal-unico")

	aesKey, hmacKey, err := pbkdf2.DerivationFunc(sharedKey, salt)

	assert.NoError(t, err)
	assert.Len(t, aesKey, pbkdf2.AESKeySize)
	assert.Len(t, hmacKey, pbkdf2.HMACKeySize)
}
