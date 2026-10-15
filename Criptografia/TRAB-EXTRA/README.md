# Relatório do Trabalho de Criptografia

Este documento descreve a implementação e os desafios encontrados no desenvolvimento do trabalho extra de Criptografia, focado em funções de hash (SHAKE128).

## Como Rodar

Certifique-se de ter a linguagem Go instalada.

No diretório raiz do projeto, execute o comando:

```bash
go mod tidy
```
Esse comando baixará todas as dependências necessárias para rodar o projeto.
Em seguida, execute:

```bash
go run main.go
```

Isso executará a função main, que os desafios implementados e exibirá os resultados no terminal.

Caso queira rodar apenas os testes para comparar os hashes, execute:
```bash
go run main.go -testes
```

## Passo a Passo da Implementação

O trabalho foi dividido em três desafios principais, todos utilizando a função de hash SHAKE128 da biblioteca `golang.org/x/crypto/sha3`.

### Função Auxiliar (`getSHAKE128`)
Foi criada uma função wrapper que recebe uma string de entrada e o tamanho desejado de saída em bytes. Ela instancia o SHAKE128, escreve a entrada e lê a quantidade de bytes solicitada para o buffer de saída.

### Desafio A: Resistência a Colisão
**Objetivo**: Encontrar dois inputs distintos que produzam o mesmo hash de 4 bytes (32 bits).
**Estratégia**:
1.  Utilizou-se um `map[string]string` para armazenar os hashes já encontrados e seus respectivos inputs originais.
2.  A busca foi paralelizada utilizando goroutines (baseado no número de CPUs disponíveis).
3.  Cada goroutine gera inputs baseados em um contador atômico (`col_seed_N`), calcula o hash e verifica se ele já existe no mapa.
4.  Se o hash já existe no mapa, uma colisão foi encontrada e o processo é encerrado.

### Desafio B: Segunda Pré-imagem
**Objetivo**: Dado um input fixo (`x1`), encontrar um input diferente (`x2`) que tenha o mesmo hash de 4 bytes.
**Estratégia**:
1.  Calculou-se o hash alvo a partir do input com o nome do aluno.
2.  Iniciou-se uma busca por força bruta paralela gerando inputs do tipo `nonce_N`.
3.  Comparou-se o hash de cada `nonce` com o hash alvo até encontrar uma correspondência.

### Desafio C: Pré-imagem (Prefixo)
**Objetivo**: Encontrar um input cujo hash corresponda a um prefixo hexadecimal específico.
**Estratégia**:
1.  O alvo foi fornecido como uma string hexadecimal.
2.  A busca paralela gera inputs `pwd_N` e calcula um hash de 4 bytes.
3.  Verifica-se se a representação hexadecimal do hash gerado é o mesmo que o alvo fornecido.

## Dificuldades Encontradas

1.  **Concorrência e Acesso a Mapas**:
    No Desafio A, a principal dificuldade foi gerenciar o acesso concorrente ao mapa de hashes vistos (`seen`). Em Go, mapas não são *thread-safe* para escrita simultânea, o que causaria "panic" em tempo de execução.
    *Solução*: Foi introduzido um `sync.Mutex` para bloquear a execução (`Lock`/`Unlock`) durante a leitura e escrita no mapa, garantindo a integridade dos dados.

2.  **Sincronização e Cancelamento**:
    Coordenar o encerramento de todas as goroutines assim que uma solução fosse encontrada foi um ponto de atenção. Sem um mecanismo de parada, as outras threads continuariam consumindo CPU desnecessariamente após o sucesso de uma delas.
    *Solução*: Utilizou-se o `context.WithCancel`. Assim que uma goroutine encontra a solução, ela chama a função `cancel()`, sinalizando através do canal `ctx.Done()` que as outras goroutines devem parar.

3.  **Geração de Entropia e Contadores**:
    Garantir que cada thread testasse valores únicos sem sobreposição de trabalho.
    *Solução*: Uso de `sync/atomic` para incrementar um contador global (`i`) de forma segura entre as threads, garantindo que cada iteração teste um `seed` ou `nonce` único.