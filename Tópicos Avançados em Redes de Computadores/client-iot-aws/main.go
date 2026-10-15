package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

var broadcast = make(chan []byte, 1024)

var upgrader = websocket.Upgrader{}

func main() {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go connectToAWSIoT()

	go handleMessages()

	fmt.Println("Servidor rodando em http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro no upgrade do WebSocket:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	log.Println("Novo cliente WebSocket conectado")

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Cliente desconectado: %v", err)
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Erro ao enviar mensagem para o cliente: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func connectToAWSIoT() {
	awsEndpoint := "a65asywtuvml2-ats.iot.us-east-1.amazonaws.com"
	mqttTopic := "esp32/dht11"
	clientID := "go-backend-server"

	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:8883", awsEndpoint))
	opts.SetClientID(clientID)
	opts.SetTLSConfig(newTLSConfig())

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		log.Printf("Mensagem recebida do tópico [%s]: %s\n", msg.Topic(), msg.Payload())
		broadcast <- msg.Payload()
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao conectar no AWS IoT: %s", token.Error())
	}

	log.Println("Conectado com sucesso ao AWS IoT!")

	if token := client.Subscribe(mqttTopic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao se inscrever no tópico: %s", token.Error())
	}

	log.Printf("Inscrito com sucesso no tópico: %s", mqttTopic)
}

func newTLSConfig() *tls.Config {
	certpool := x509.NewCertPool()
	caCertPath := filepath.Join("certs", "AmazonRootCA1.pem")

	pemCerts, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Erro ao ler o certificado CA raiz (%s): %v", caCertPath, err)
	}
	certpool.AppendCertsFromPEM(pemCerts)

	certPath := filepath.Join("certs", "certificate.pem.crt")
	keyPath := filepath.Join("certs", "private.pem.key")
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("Erro ao carregar o par de chaves do cliente: %v", err)
	}

	return &tls.Config{
		RootCAs:            certpool,
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}
}
