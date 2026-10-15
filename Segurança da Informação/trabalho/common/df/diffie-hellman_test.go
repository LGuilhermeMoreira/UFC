package df_test

import (
	"testing"
	"trabalho-seguranca/common/df"

	"github.com/stretchr/testify/assert"
)

func TestDiffieHellman(t *testing.T) {
	// Alice
	a, err := df.GenerateSecretKey(df.Prime)
	assert.NoError(t, err)

	A, err := df.GenerateSharedKeys(df.Prime, df.Generator, a)
	assert.NoError(t, err)

	//Bob
	b, err := df.GenerateSecretKey(df.Prime)
	assert.NoError(t, err)

	B, err := df.GenerateSharedKeys(df.Prime, df.Generator, b)
	assert.NoError(t, err)

	// Alice
	sharedSecretKey1, err := df.GenerateSharedSecretKey(df.Prime, B, a)
	assert.NoError(t, err)

	//Bob
	sharedSecretKey2, err := df.GenerateSharedSecretKey(df.Prime, A, b)

	assert.Equal(t, *sharedSecretKey1, *sharedSecretKey2)
}
