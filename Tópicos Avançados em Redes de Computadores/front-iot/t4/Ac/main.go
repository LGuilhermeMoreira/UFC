package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"test-fyni/t4/ex"
	"test-fyni/t4/model"
	"test-fyni/t4/mqtt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func MockSubscribe(ch chan string) {
	go func() {
		for _, v := range ex.Data {
			time.Sleep(1 * time.Second)
			ch <- v
		}
	}()
}

// Histórico dos dados
var axData []float32
var ayData []float32
var azData []float32

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.LightTheme())
	myWindow := myApp.NewWindow("Gráfico de Aceleração AX")

	ch := make(chan string)
	updateChan := make(chan func())

	graphContainer := container.NewWithoutLayout()

	// >>>>>> AQUI criamos os botões <<<<<<
	startButton := widget.NewButton("Start", func() {
		go mqtt.SubscribeMQTT("tcp://broker.hivemq.com:1883", "/home/accel", ch)
		// go MockSubscribe(ch)
	})

	clearButton := widget.NewButton("Limpar", func() {
		axData = nil
		ayData = nil
		azData = nil
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

	legendContainer := container.NewHBox(layout.NewSpacer())

	legendas := []struct {
		text  string
		color color.Color
	}{
		{"AX", color.RGBA{0, 200, 0, 255}}, // Verde
		{"AY", color.RGBA{0, 0, 255, 255}}, // Azul
		{"AZ", color.RGBA{255, 0, 0, 255}}, // Vermelho
	}

	for _, l := range legendas {
		rect := canvas.NewRectangle(l.color)
		rect.Resize(fyne.NewSize(15, 15))
		legendContainer.Add(rect)

		label := canvas.NewText(l.text, color.Black)
		label.TextSize = 14
		label.Move(fyne.NewPos(20, -2)) // Ajuste fino da posição do texto em relação ao retângulo
		legendContainer.Add(label)
		legendContainer.Add(layout.NewSpacer()) // Espaçamento entre as legendas
	}

	// >>>>>> Novo container principal <<<<<<
	mainContainer := container.NewBorder(
		container.NewVBox( // Topo
			widget.NewLabelWithStyle("Gráfico de Aceleração AX", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			buttonContainer,
			container.NewHBox(layout.NewSpacer(), legendContainer, layout.NewSpacer()), // Centraliza as legendas
		),
		nil,            // sem bottom
		nil,            // sem left
		nil,            // sem right
		graphContainer, // centro é o gráfico
	)

	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(800, 600))

	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(800, 600))
	// MockSubscribe(ch)

	// Goroutine que escuta atualizações para a UI
	go func() {
		for update := range updateChan {
			update() // Executa na thread principal
		}
	}()

	go func() {
		for msg := range ch {
			var data model.SensorData
			if err := json.Unmarshal([]byte(msg), &data); err != nil {
				log.Println("Erro ao fazer unmarshal:", err)
				continue
			}

			axData = append(axData, data.AX)
			ayData = append(ayData, data.AY)
			azData = append(azData, data.AZ)

			if len(axData) > 100 {
				axData = axData[1:]
				ayData = ayData[1:]
				azData = azData[1:]
			}

			updateChan <- func() {
				width := float32(800)
				height := float32(600)
				margin := float32(50)
				graphWidth := width - 2*margin
				graphHeight := height - 2*margin
				step := graphWidth / float32(len(axData))

				maxAbsValue := float32(1.0)
				escala := graphHeight / (2 * maxAbsValue)

				graphContainer.Objects = nil

				// Eixos
				eixoY := canvas.NewLine(color.Black)
				eixoY.StrokeWidth = 2
				eixoY.Position1 = fyne.NewPos(margin, margin)
				eixoY.Position2 = fyne.NewPos(margin, height-margin)
				graphContainer.Add(eixoY)

				eixoX := canvas.NewLine(color.Black)
				eixoX.StrokeWidth = 2
				eixoX.Position1 = fyne.NewPos(margin, height/2)
				eixoX.Position2 = fyne.NewPos(width-margin, height/2)
				graphContainer.Add(eixoX)

				// Grade e rótulos Y
				numGrades := 4
				for i := -numGrades; i <= numGrades; i++ {
					value := float32(i) * maxAbsValue / float32(numGrades)
					y := height/2 - value*escala

					grade := canvas.NewLine(color.RGBA{200, 200, 200, 100})
					grade.Position1 = fyne.NewPos(margin, y)
					grade.Position2 = fyne.NewPos(width-margin, y)
					graphContainer.Add(grade)

					label := canvas.NewText(fmt.Sprintf("%.2f", value), color.Black)
					label.TextSize = 12
					label.Move(fyne.NewPos(5, y-6))
					graphContainer.Add(label)
				}

				// Rótulos
				labelY := canvas.NewText("Aceleração", color.Black)
				labelY.TextSize = 14
				labelY.Move(fyne.NewPos(5, height/2))
				graphContainer.Add(labelY)

				labelX := canvas.NewText("Tempo (amostras)", color.Black)
				labelX.TextSize = 14
				labelX.Move(fyne.NewPos(width/2-50, height-30))
				graphContainer.Add(labelX)

				// Desenho AX
				for i := 1; i < len(axData); i++ {
					x1 := margin + step*float32(i-1)
					x2 := margin + step*float32(i)
					y1 := height/2 - axData[i-1]*escala
					y2 := height/2 - axData[i]*escala

					line := canvas.NewLine(color.RGBA{0, 200, 0, 255})
					line.StrokeWidth = 2
					line.Position1 = fyne.NewPos(x1, y1)
					line.Position2 = fyne.NewPos(x2, y2)
					graphContainer.Add(line)
				}

				// AY
				for i := 1; i < len(ayData); i++ {
					x1 := margin + step*float32(i-1)
					x2 := margin + step*float32(i)
					y1 := height/2 - ayData[i-1]*escala
					y2 := height/2 - ayData[i]*escala

					line := canvas.NewLine(color.RGBA{0, 0, 255, 255})
					line.StrokeWidth = 2
					line.Position1 = fyne.NewPos(x1, y1)
					line.Position2 = fyne.NewPos(x2, y2)
					graphContainer.Add(line)
				}

				// AZ
				for i := 1; i < len(azData); i++ {
					x1 := margin + step*float32(i-1)
					x2 := margin + step*float32(i)
					y1 := height/2 - azData[i-1]*escala
					y2 := height/2 - azData[i]*escala

					line := canvas.NewLine(color.RGBA{255, 0, 0, 255})
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
