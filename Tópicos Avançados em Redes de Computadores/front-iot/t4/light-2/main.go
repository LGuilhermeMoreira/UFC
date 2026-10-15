package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"test-fyni/t4/mqtt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// Altere para o caminho correto do seu módulo
)

// LightSensor representa os dados recebidos do MQTT
type LightSensor struct {
	Light int `json:"light"`
}

func main() {
	// Cria o canal para receber mensagens MQTT
	messageChan := make(chan string)

	// Inicia a escuta MQTT
	go mqtt.SubscribeMQTT("tcp://broker.hivemq.com:1883", "/home/ldr", messageChan)

	// Cria a aplicação Fyne
	myApp := app.New()
	myApp.Settings().SetTheme(theme.LightTheme())
	myWindow := myApp.NewWindow("Sensor de Luz - 2")

	// Cria um label que será atualizado com os dados do MQTT
	statusLabel := widget.NewLabel("Aguardando dados do sensor...")

	// Cria um círculo para a animação
	circle := canvas.NewCircle(color.Gray{})
	circle.Resize(fyne.NewSize(50, 50))

	// Coloca o label e o círculo em um container
	content := container.NewVBox(
		widget.NewLabel("Status da Luz:"),
		statusLabel,
		circle,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 200)) // Aumenta um pouco a altura para o círculo
	myWindow.Show()

	// Atualiza o label e a cor do círculo sempre que novas mensagens chegarem
	go func() {
		for msg := range messageChan {
			var sensor LightSensor
			err := json.Unmarshal([]byte(msg), &sensor)
			if err != nil {
				log.Println("Erro ao fazer Unmarshal:", err)
				continue
			}

			status := ""
			if sensor.Light == 0 {
				status = "Ambiente com luz suficiente."
				circle.FillColor = color.RGBA{G: 255, A: 255} // Verde
			} else {
				status = "Ambiente com pouca luz."
				circle.FillColor = color.RGBA{R: 255, A: 255} // Vermelho
			}

			statusLabel.SetText(fmt.Sprintf("%s\n(recebido às %s)", status, time.Now().Format("15:04:05")))
			canvas.Refresh(circle) // Força a atualização do círculo
			content.Refresh()
		}
	}()

	myApp.Run()
}
