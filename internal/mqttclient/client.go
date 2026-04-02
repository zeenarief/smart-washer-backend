package mqttclient

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// InitMQTT menginisialisasi koneksi ke broker Mosquitto
func InitMQTT() mqtt.Client {
	broker := os.Getenv("MQTT_BROKER")
	user := os.Getenv("MQTT_USER")
	pass := os.Getenv("MQTT_PASS")
	clientID := os.Getenv("MQTT_CLIENT_ID")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	// Jika Mosquitto Anda menggunakan username/password nantinya
	if user != "" && pass != "" {
		opts.SetUsername(user)
		opts.SetPassword(pass)
	}

	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	// Callback ketika koneksi terputus
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("Koneksi MQTT terputus: %v", err)
	}

	// Callback ketika berhasil terkoneksi
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("Koneksi MQTT ke broker berhasil!")
		// Nanti kita akan tambahkan logika Subscribe di sini
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Gagal terhubung ke MQTT Broker: %v", token.Error())
	}

	return client
}

// PublishCommand adalah helper untuk mengirim perintah ke mesin cuci
func PublishCommand(client mqtt.Client, macAddress string, payload string) error {
	topic := fmt.Sprintf("mesincuci/%s/command", macAddress)
	token := client.Publish(topic, 1, false, payload)
	token.Wait()
	return token.Error()
}
