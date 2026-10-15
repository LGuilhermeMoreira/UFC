# Sistema Solar 3D

## Sobre

Este é um projeto desenvolvido para a cadeira de Computação Gráfica 2026.1 da UFC - Quixadá.

O projeto simula um sistema solar em 3D com planetas texturizados, orbitas elipticas, luas, aneis de Saturno, fundo estrelado e camera orbital interativa. A aplicacao permite pausar a simulacao, alterar a velocidade, focar planetas e alternar elementos visuais como labels, orbitas e grade.

## Stack

- Go
- raylib-go
- Raylib
- Makefile

## Dependencias

Dependencias Go declaradas em `go.mod`:

- `github.com/gen2brain/raylib-go/raylib v0.60.0`
- `github.com/ebitengine/purego v0.10.0` indireta
- `github.com/jupiterrider/ffi v0.7.0` indireta
- `golang.org/x/exp` indireta

Tambem e necessario ter o Go instalado. Dependendo do sistema operacional, a raylib-go pode exigir bibliotecas nativas de compilacao e video/audio instaladas no ambiente.

## Como rodar

Clone o repositorio e entre na pasta do projeto:

```bash
git clone <url-do-repositorio>
cd sistema-solar
```

Baixe as dependencias:

```bash
go mod download
```

Execute diretamente:

```bash
go run .
```

Ou compile usando o Makefile:

```bash
make build
./sistema-solar
```

Para remover o binario gerado:

```bash
make clean
```

## Controles

- `Espaco`: pausar ou continuar a simulacao
- `Mouse`: arrastar para rotacionar a camera
- `Scroll`: aproximar ou afastar a camera
- `Clique em um planeta`: focar o planeta selecionado
- `W`, `A`, `S`, `D`: mover o alvo da camera no plano horizontal
- `Q`, `E`: mover o alvo da camera no eixo vertical
- `Seta para cima` e `Seta para baixo`: aumentar ou reduzir a velocidade
- `Z`: resetar a velocidade
- `L`: mostrar ou ocultar labels
- `O`: mostrar ou ocultar orbitas
- `G`: mostrar ou ocultar grade
- `C`: lancar um cometa contra a Terra
- `R`: resetar camera e velocidade

## Estrutura

- `main.go`: inicializacao da janela, estado da simulacao, entrada do usuario e loop principal
- `models.go`: estruturas usadas por planetas, luas, camera e estrelas
- `helper.go`: funcoes auxiliares para camera, desenho, orbitas, selecao e carregamento de assets
- `assets/`: texturas dos planetas, Sol e aneis de Saturno
- `Makefile`: comandos de build e limpeza
