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
	"fyne.io/fyne/v2/theme"
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
			"aZ": "-0.52 g",
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
var ayData []float64
var azData []float64

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.LightTheme())
	myWindow := myApp.NewWindow("Gráfico de Aceleração AX")

	ch := make(chan string)
	updateChan := make(chan func())

	graphContainer := container.NewWithoutLayout()

	// >>>>>> AQUI criamos os botões <<<<<<
	startButton := widget.NewButton("Start", func() {
		go Subscribe(ch)
	})

	clearButton := widget.NewButton("Limpar", func() {
		axData = nil
		updateChan <- func() {
			graphContainer.Objects = nil
			graphContainer.Refresh()
		}
	})

	// Agrupamos os botões horizontalmente
	buttonContainer := container.NewHBox(
		layout.NewSpacer(),
		startButton,
		layout.NewSpacer(),
		clearButton,
		layout.NewSpacer(),
	)

	// >>>>>> Novo container principal <<<<<<
	mainContainer := container.NewBorder(
		container.NewVBox( // Topo
			widget.NewLabelWithStyle("Gráfico de Aceleração AX", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			buttonContainer,
		),
		nil,            // sem bottom
		nil,            // sem left
		nil,            // sem right
		graphContainer, // centro é o gráfico
	)

	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(800, 600))
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

			ayValue, err := strconv.ParseFloat(strings.TrimSuffix(data.AY, " g"), 64)
			if err != nil {
				log.Println("Erro ao converter aY:", err)
				continue
			}

			azValue, err := strconv.ParseFloat(strings.TrimSuffix(data.AZ, " g"), 64)
			if err != nil {
				log.Println("Erro ao converter aZ:", err)
				continue
			}

			axData = append(axData, axValue)
			ayData = append(ayData, ayValue)
			azData = append(azData, azValue)

			if len(axData) > 100 {
				axData = axData[1:]
				ayData = ayData[1:]
				azData = azData[1:]
			}

			updateChan <- func() {

				width := float32(800)
				height := float32(600)
				margin := float32(50) // Margem para eixo e rótulos
				graphWidth := width - 2*margin
				graphHeight := height - 2*margin
				step := graphWidth / float32(len(axData))

				graphContainer.Objects = nil // Limpa

				// --- Desenha os eixos ---
				eixoX := canvas.NewLine(color.Black)
				eixoX.StrokeWidth = 2
				eixoX.Position1 = fyne.NewPos(margin, height-margin)
				eixoX.Position2 = fyne.NewPos(width-margin, height-margin)
				graphContainer.Add(eixoX)

				eixoY := canvas.NewLine(color.Black)
				eixoY.StrokeWidth = 2
				eixoY.Position1 = fyne.NewPos(margin, margin)
				eixoY.Position2 = fyne.NewPos(margin, height-margin)
				graphContainer.Add(eixoY)

				// --- Adiciona Título ---
				titulo := canvas.NewText("Gráfico de Aceleração no Eixo X,Y,Z", color.Black)
				titulo.Alignment = fyne.TextAlignCenter
				titulo.TextStyle = fyne.TextStyle{Bold: true}
				titulo.TextSize = 24
				titulo.Move(fyne.NewPos(width/2-150, 10))
				graphContainer.Add(titulo)

				// Adiciona legenda colorida
				legendX := float32(width - margin - 120)
				legendY := float32(50)

				legendas := []struct {
					text string
					col  color.Color
				}{
					{"AX", color.RGBA{0, 200, 0, 255}}, // Verde
					{"AY", color.RGBA{0, 0, 255, 255}}, // Azul
					{"AZ", color.RGBA{255, 0, 0, 255}}, // Vermelho
				}

				for i, l := range legendas {
					rect := canvas.NewRectangle(l.col)
					rect.Resize(fyne.NewSize(15, 15))
					rect.Move(fyne.NewPos(legendX, legendY+float32(i*25)))
					graphContainer.Add(rect)

					label := canvas.NewText(l.text, color.Black)
					label.TextSize = 14
					label.Move(fyne.NewPos(legendX+20, legendY+float32(i*25)-2))
					graphContainer.Add(label)
				}

				numGrades := 5
				for i := 0; i <= numGrades; i++ {
					y := margin + graphHeight/float32(numGrades)*float32(i)
					grade := canvas.NewLine(color.RGBA{150, 150, 150, 100}) // Cinza claro
					grade.Position1 = fyne.NewPos(margin, height-y)
					grade.Position2 = fyne.NewPos(width-margin, height-y)
					graphContainer.Add(grade)
				}

				// --- Rótulos dos eixos ---
				labelY := canvas.NewText("Aceleração (g)", color.Black)
				labelY.TextSize = 14
				labelY.Move(fyne.NewPos(5, height/2))
				graphContainer.Add(labelY)

				labelX := canvas.NewText("Tempo (amostras)", color.Black)
				labelX.TextSize = 14
				labelX.Move(fyne.NewPos(width/2-50, height-30))
				graphContainer.Add(labelX)

				// --- Desenha o gráfico dos dados ---
				for i := 1; i < len(axData); i++ {
					x1 := margin + step*float32(i-1)
					y1 := height - margin - float32(axData[i-1]*200)
					x2 := margin + step*float32(i)
					y2 := height - margin - float32(axData[i]*200)

					line := canvas.NewLine(color.RGBA{0, 200, 0, 255})
					line.StrokeWidth = 2
					line.Position1 = fyne.NewPos(x1, y1)
					line.Position2 = fyne.NewPos(x2, y2)
					graphContainer.Add(line)
				}

				for i := 1; i < len(ayData); i++ {
					x1 := margin + step*float32(i-1)
					y1 := height - margin - float32(ayData[i-1]*200)
					x2 := margin + step*float32(i)
					y2 := height - margin - float32(ayData[i]*200)

					line := canvas.NewLine(color.RGBA{0, 0, 255, 255}) // Azul
					line.StrokeWidth = 2
					line.Position1 = fyne.NewPos(x1, y1)
					line.Position2 = fyne.NewPos(x2, y2)
					graphContainer.Add(line)
				}

				for i := 1; i < len(azData); i++ {
					x1 := margin + step*float32(i-1)
					y1 := height - margin - float32(azData[i-1]*200)
					x2 := margin + step*float32(i)
					y2 := height - margin - float32(azData[i]*200)

					line := canvas.NewLine(color.RGBA{255, 0, 0, 255}) // Vermelho
					line.StrokeWidth = 2
					line.Position1 = fyne.NewPos(x1, y1)
					line.Position2 = fyne.NewPos(x2, y2)
					graphContainer.Add(line)
				}

				graphContainer.Refresh()

			}
		}
	}()

	myWindow.ShowAndRun()
}
