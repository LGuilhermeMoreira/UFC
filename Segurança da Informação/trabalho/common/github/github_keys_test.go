package github_test

import (
	"sync"
	"testing"
	"trabalho-seguranca/common/github"

	"github.com/stretchr/testify/assert"
)

func TestGithubKeys(t *testing.T) {
	var chaveClient, chaveServidor string
	var erroClient, erroServidor error
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		chaveClient, erroClient = github.LoadClientECDSAPublicKey()
	}()

	go func() {
		defer wg.Done()
		chaveServidor, erroServidor = github.LoadServerECDSAPublicKey()
	}()

	wg.Wait()

	assert.NoError(t, erroClient)
	assert.NoError(t, erroServidor)

	assert.NotEmpty(t, chaveClient)
	assert.NotEmpty(t, chaveServidor)

	assert.NotEqual(t, "/n", chaveClient)
	assert.NotEqual(t, "/n", chaveServidor)
}
