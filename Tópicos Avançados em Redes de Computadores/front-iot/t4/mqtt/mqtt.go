package mqtt

import (
	"fmt"
	"log"

	mqttPkg "github.com/eclipse/paho.mqtt.golang"
)

func SubscribeMQTT(broker string, topic string, ch chan string) {
	opts := mqttPkg.NewClientOptions().AddBroker(broker)
	opts.SetClientID("fyne-subscriber")
	opts.SetDefaultPublishHandler(func(client mqttPkg.Client, msg mqttPkg.Message) {
		fmt.Printf("\nmessage: %s\n", string(msg.Payload()))
		ch <- string(msg.Payload())
	})

	client := mqttPkg.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Erro ao conectar ao MQTT:", token.Error())
		return
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Println("Erro ao se inscrever no tópico:", token.Error())
		return
	}

	log.Println("Conectado e inscrito no tópico MQTT:", topic)
}
