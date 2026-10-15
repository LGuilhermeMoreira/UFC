package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/sha3"
)

func getSHAKE128(input string, outputSize int) []byte {
	// crio um buffer para armazenar o resultado
	hash := make([]byte, outputSize)
	// crio um "shake"
	d := sha3.NewShake128()
	// shake processa o input
	d.Write([]byte(input))
	// copio o resultado do shake no buffer
	d.Read(hash)
	// retorno o buffer
	return hash
}

// RESISTÊNCIA A COLISÃO
func DesafioA() {
	fmt.Println("--- Desafio A ---")
	// crio um map
	seen := make(map[string]string)
	// crio o mutex
	var mu sync.Mutex

	var i int64 = 0
	// crio um context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// crio um canal
	done := make(chan bool)

	// cria um thread pool
	for w := 0; w < runtime.NumCPU(); w++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// pega um i para cada seed
					valI := atomic.AddInt64(&i, 1)
					// gera o input do shake
					input := fmt.Sprintf("col_seed_%d", valI)
					// pega a saida do shake e faz um encode para string
					hash := hex.EncodeToString(getSHAKE128(input, 4))

					// loca a thread, map não é thread safe
					mu.Lock()
					if original, exists := seen[hash]; exists {
						// se der bom:

						// printo o resultado
						fmt.Printf("Colisão encontrada!\nInput 1: %s\nInput 2: %s\nHash: %s\n\n", original, input, hash)
						// cancelo o context
						cancel()
						// unlock na thared
						mu.Unlock()
						// sinalizo as thareads que deu bom
						done <- true
						// saio da go func
						return
					}
					seen[hash] = input
					// faz o unlock da thread
					mu.Unlock()
				}
			}
		}()
	}
	// sincronizar as threads
	<-done
}

// SEGUNDA PRÉ-IMAGEM
func DesafioB(nomeCompleto string) {
	fmt.Println("--- Desafio B ---")
	x1 := "Aluno: " + nomeCompleto
	// gero o hash alvo
	targetHash := hex.EncodeToString(getSHAKE128(x1, 4))
	fmt.Printf("Alvo (x1): %s | Hash: %s\n", x1, targetHash)

	// mesmo passo do desafio A
	var i int64 = 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan bool)

	for w := 0; w < runtime.NumCPU(); w++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					valI := atomic.AddInt64(&i, 1)
					x2 := fmt.Sprintf("nonce_%d", valI)
					// caso forema  mesma saida eu pulo a interação
					if x2 == x1 {
						continue
					}

					// vejo se achei o hash alvo
					if hex.EncodeToString(getSHAKE128(x2, 4)) == targetHash {
						fmt.Printf("Segunda pré-imagem encontrada!\nInput: %s\nHash: %s\n\n", x2, targetHash)
						cancel()
						done <- true
						return
					}
				}
			}
		}()
	}
	<-done
}

// PRÉ-IMAGEM 34 BITS
func DesafioC(targetHex string, targetBytes int) {
	fmt.Printf("--- Desafio C ---\n")

	var i int64 = 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan bool)

	for w := 0; w < runtime.NumCPU(); w++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					valI := atomic.AddInt64(&i, 1)
					input := fmt.Sprintf("pwd_%d", valI)

					hash := getSHAKE128(input, targetBytes)

					hashHex := strings.ToUpper(hex.EncodeToString(hash))
					if hashHex == targetHex {
						fmt.Printf("Pré-imagem encontrada!\nSenha: %s\nHash: %s\n\n", input, hashHex)
						cancel()
						done <- true
						return
					}
				}
			}
		}()
	}
	<-done
}

func runTests() {
	fmt.Println("Test 1")
	fmt.Printf("1: %v\n2: %v\nhash1: %v\nhash2: %v\n",
		"col_seed_19075",
		"col_seed_67494",
		hex.EncodeToString(getSHAKE128("col_seed_19075", 4)),
		hex.EncodeToString(getSHAKE128("col_seed_67494", 4)))

	fmt.Println("Test 2")
	fmt.Printf("1: %v\n2: %v\nhash1: %v\nhash2: %v\n",
		"Aluno: Luiz Guilherme Moreira Leite",
		"nonce_575534485",
		hex.EncodeToString(getSHAKE128("Aluno: Luiz Guilherme Moreira Leite", 4)),
		hex.EncodeToString(getSHAKE128("nonce_575534485", 4)))

	fmt.Println("Test 3")
	fmt.Printf("1: %v\nhash1: %v\nhash2: %v\n",
		"pwd_731269163",
		"17675FC0",
		hex.EncodeToString(getSHAKE128("pwd_731269163", 4)))
}

func main() {
	testes := flag.Bool("testes", false, "Roda os testes")
	flag.Parse()

	if *testes {
		runTests()
	} else {
		t1 := time.Now()
		DesafioA()
		DesafioB("Luiz Guilherme Moreira Leite")
		DesafioC("17675FC0", 4)
		t2 := time.Now()
		fmt.Printf("Tempo total de execução: %v\n", t2.Sub(t1))
	}

}

/*
--- Desafio A ---
Colisão encontrada!
Input 1: col_seed_19075
Input 2: col_seed_67494
Hash: c903f0c2

--- Desafio B ---
Alvo (x1): Aluno: Luiz Guilherme Moreira Leite | Hash: 9627ddc0
Segunda pré-imagem encontrada!
Input: nonce_575534485
Hash: 9627ddc0

--- Desafio C ---
Pré-imagem encontrada!
Senha: pwd_731269163
Hash: 17675FC03D
*/
