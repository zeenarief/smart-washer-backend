package mqttclient

import (
	"encoding/json"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeenarief/smart-washer-backend/internal/ws"
)

// SetupSubscriber mendengarkan pesan status dari ESP32
func SetupSubscriber(client mqtt.Client, hub *ws.Hub) {
	topic := "mesincuci/+/status"
	token := client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
		log.Printf("Terima status dari ESP32: %s", string(m.Payload()))

		var status map[string]interface{}
		if err := json.Unmarshal(m.Payload(), &status); err != nil {
			log.Printf("Gagal parse MQTT status: %v", err)
			return
		}

		// EKSTRAKSI MAC ADDRESS DARI TOPIK
		// Topik: mesincuci/[MAC]/status
		topicParts := strings.Split(m.Topic(), "/")
		if len(topicParts) >= 3 {
			// Masukkan MAC ke dalam payload agar Flutter tahu ini milik siapa
			status["mac"] = topicParts[1]
		}

		// Teruskan JSON (yang kini berisi data wash, spin, dan mac) ke Flutter
		hub.BroadcastStatus(status)
	})

	token.Wait()
	log.Printf("Subscribed ke topik: %s", topic)
}
