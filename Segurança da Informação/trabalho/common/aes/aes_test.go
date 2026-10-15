package aes_test

import (
	"crypto/sha256"
	"testing"
	"trabalho-seguranca/common/aes"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {
	message := []byte("Fala meu cumparça!")
	aeskey := sha256.Sum256([]byte("123321222"))
	key := aeskey[:]
	ct, iv, err := aes.Encrypt(key, message)
	assert.NoError(t, err)
	assert.NotEmpty(t, ct)
	assert.NotEmpty(t, iv)

	pt, err := aes.Decrypt(key, iv, ct)
	assert.NoError(t, err)
	assert.Equal(t, string(message), string(pt))
}
