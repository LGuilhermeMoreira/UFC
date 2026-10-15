package main

import (
	"encoding/json"
	"fmt"
	"log"
	"test-fyni/t4/mqtt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	myWindow := myApp.NewWindow("Sensor de Luz")

	// Cria um label que será atualizado com os dados do MQTT
	statusLabel := widget.NewLabel("Aguardando dados do sensor...")

	// Coloca o label em um container
	content := container.NewVBox(
		widget.NewLabel("Status da Luz:"),
		statusLabel,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 150))
	myWindow.Show()

	// Atualiza o label sempre que novas mensagens chegarem
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
			} else {
				status = "Ambiente com pouca luz."
			}

			statusLabel.SetText(fmt.Sprintf("%s\n(recebido às %s)", status, time.Now().Format("15:04:05")))

			content.Refresh()
		}
	}()

	myApp.Run()
}
