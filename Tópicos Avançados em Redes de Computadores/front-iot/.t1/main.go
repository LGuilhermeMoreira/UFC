package main

import (
	"encoding/json"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type SensorData struct {
	AX string `json:"aX"`
	AY string `json:"aY"`
	AZ string `json:"aZ"`
	GX string `json:"gX"`
	GY string `json:"gY"`
	GZ string `json:"gZ"`
}

func Subscribe(ch chan string) {
	arr := []string{
		`{
			"aX": "0.12 g",
			"aY": "0.08 g",
			"aZ": "0.50 g",
			"gX": "1.2 dps",
			"gY": "1.8 dps",
			"gZ": "0.5 dps"
		}`,
		`{
			"aX": "0.20 g",
			"aY": "0.10 g",
			"aZ": "0.45 g",
			"gX": "1.5 dps",
			"gY": "2.1 dps",
			"gZ": "0.7 dps"
		}`,
		`{
			"aX": "0.15 g",
			"aY": "0.09 g",
			"aZ": "0.48 g",
			"gX": "1.3 dps",
			"gY": "1.9 dps",
			"gZ": "0.6 dps"
		}`,
		`{
			"aX": "0.18 g",
			"aY": "0.11 g",
			"aZ": "0.52 g",
			"gX": "1.6 dps",
			"gY": "2.2 dps",
			"gZ": "0.8 dps"
		}`,
		`{
			"aX": "0.13 g",
			"aY": "0.07 g",
			"aZ": "0.49 g",
			"gX": "1.1 dps",
			"gY": "1.7 dps",
			"gZ": "0.4 dps"
		}`,
		`{
			"aX": "0.22 g",
			"aY": "0.12 g",
			"aZ": "0.47 g",
			"gX": "1.7 dps",
			"gY": "2.3 dps",
			"gZ": "0.9 dps"
		}`,
		`{
			"aX": "0.16 g",
			"aY": "0.085 g",
			"aZ": "0.51 g",
			"gX": "1.4 dps",
			"gY": "2.0 dps",
			"gZ": "0.55 dps"
		}`,
		`{
			"aX": "0.19 g",
			"aY": "0.105 g",
			"aZ": "0.46 g",
			"gX": "1.55 dps",
			"gY": "2.15 dps",
			"gZ": "0.75 dps"
		}`,
		`{
			"aX": "0.14 g",
			"aY": "0.095 g",
			"aZ": "0.495 g",
			"gX": "1.25 dps",
			"gY": "1.85 dps",
			"gZ": "0.65 dps"
		}`,
		`{
			"aX": "0.21 g",
			"aY": "0.115 g",
			"aZ": "0.455 g",
			"gX": "1.65 dps",
			"gY": "2.25 dps",
			"gZ": "0.85 dps"
		}`,
		`{
			"aX": "0.17 g",
			"aY": "0.09 g",
			"aZ": "0.505 g",
			"gX": "1.35 dps",
			"gY": "1.95 dps",
			"gZ": "0.5 dps"
		}`,
		`{
			"aX": "0.125 g",
			"aY": "0.085 g",
			"aZ": "0.485 g",
			"gX": "1.15 dps",
			"gY": "1.75 dps",
			"gZ": "0.6 dps"
		}`,
		`{
			"aX": "0.205 g",
			"aY": "0.105 g",
			"aZ": "0.465 g",
			"gX": "1.55 dps",
			"gY": "2.15 dps",
			"gZ": "0.7 dps"
		}`,
		`{
			"aX": "0.155 g",
			"aY": "0.095 g",
			"aZ": "0.49 g",
			"gX": "1.3 dps",
			"gY": "1.9 dps",
			"gZ": "0.65 dps"
		}`,
		`{
			"aX": "0.185 g",
			"aY": "0.11 g",
			"aZ": "0.515 g",
			"gX": "1.6 dps",
			"gY": "2.2 dps",
			"gZ": "0.8 dps"
		}`,
	}

	go func() {
		for _, v := range arr {
			print("mandou")
			time.Sleep(1 * time.Second)
			ch <- v
		}
	}()
}

// Histórico dos dados
var axData []float64

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Gráfico de Aceleração AX")
	ch := make(chan string)
	updateChan := make(chan func())

	// Cria o container para o gráfico
	graphContainer := container.NewWithoutLayout()

	// Cria o container principal com BorderLayout e adiciona os elementos nas regiões corretas
	mainContainer := container.New(layout.NewBorderLayout(widget.NewLabel("Aceleração (g)"), nil, widget.NewLabel("Tempo (amostras)"), nil),
		graphContainer, // Adiciona o graphContainer na região central
	)
	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(800, 600))

	// Start mock MQTT
	Subscribe(ch)

	// Goroutine que escuta atualizações para a UI
	go func() {
		for update := range updateChan {
			update() // Executa na thread principal
		}
	}()

	go func() {
		for msg := range ch {
			var data SensorData
			if err := json.Unmarshal([]byte(msg), &data); err != nil {
				log.Println("Erro ao fazer unmarshal:", err)
				continue
			}

			axValue, err := strconv.ParseFloat(strings.TrimSuffix(data.AX, " g"), 64)
			if err != nil {
				log.Println("Erro ao converter aX:", err)
				continue
			}

			axData = append(axData, axValue)
			if len(axData) > 100 {
				axData = axData[1:]
			}

			updateChan <- func() {
				width := float32(800)
				height := float32(600)
				step := width / float32(len(axData))

				graphContainer.Objects = nil // Limpa os objetos do gráfico

				// Desenha os eixos
				eixoX := canvas.NewLine(color.Black)
				eixoX.Position1 = fyne.NewPos(0, height/2)
				eixoX.Position2 = fyne.NewPos(width, height/2)
				graphContainer.Add(eixoX)

				eixoY := canvas.NewLine(color.Black)
				eixoY.Position1 = fyne.NewPos(0, 0)
				eixoY.Position2 = fyne.NewPos(0, height)
				graphContainer.Add(eixoY)

				// Desenha a grade (linhas horizontais simples)
				numGrades := 5
				stepY := height / float32(numGrades)
				for i := 1; i < numGrades; i++ {
					gradeLine := canvas.NewLine(color.RGBA{100, 100, 100, 50})
					y := stepY * float32(i)
					gradeLine.Position1 = fyne.NewPos(0, y)
					gradeLine.Position2 = fyne.NewPos(width, y)
					graphContainer.Add(gradeLine)
				}

				// Desenha a linha do gráfico
				for i := 1; i < len(axData); i++ {
					x1 := step * float32(i-1)
					y1 := height/2 - float32(axData[i-1]*200)
					x2 := step * float32(i)
					y2 := height/2 - float32(axData[i]*200)

					line := canvas.NewLine(color.RGBA{0, 255, 0, 255})
					line.StrokeWidth = 2
					line.Position1 = fyne.NewPos(x1, y1)
					line.Position2 = fyne.NewPos(x2, y2)
					graphContainer.Add(line)
				}

				graphContainer.Refresh() // Atualiza o container do gráfico
			}
		}
	}()

	myWindow.ShowAndRun()
}
