package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Calculate(hmacKey, iv, ciphertext []byte) []byte {
	// calcula o hmac
	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(iv)
	mac.Write(ciphertext)
	return mac.Sum(nil)
}

func Verify(hmacKey, receivedHmac, iv, ciphertext []byte) bool {
	// calcula o novo hmac
	expectedMac := hmac.New(sha256.New, hmacKey)
	expectedMac.Write(iv)
	expectedMac.Write(ciphertext)
	expectedHmac := expectedMac.Sum(nil)

	// compara em tempo constante.
	return hmac.Equal(receivedHmac, expectedHmac)
}
