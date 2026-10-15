package hmac_test

import (
	"crypto/sha256"
	"testing"
	"trabalho-seguranca/common/hmac"

	"github.com/stretchr/testify/assert"
)

func TestHAMC(t *testing.T) {
	hmacKey := sha256.Sum256([]byte("OLA MUNDO"))
	key := hmacKey[:]
	iv := []byte("123321456789010123")
	ct := []byte("dkfmsdlfmsdkfs")
	hmacTag := hmac.Calculate(key, iv, ct)
	assert.NotEmpty(t, hmacTag)
	ok := hmac.Verify(key, hmacTag, iv, ct)
	assert.True(t, ok)
}
