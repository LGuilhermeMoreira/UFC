package genkey_test

import (
	"testing"
	genkey "trabalho-seguranca/common/gen_key"

	"github.com/stretchr/testify/assert"
)

func TestGenKeys(t *testing.T) {
	chavePrivada, chavePublica, err := genkey.GenKeys("file", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, chavePrivada)
	assert.NotEmpty(t, chavePublica)
}
