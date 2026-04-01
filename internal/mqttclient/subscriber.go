package mqttclient

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeenarief/smart-washer-backend/internal/ws"
)

// SetupSubscriber mendengarkan pesan status dari ESP32
func SetupSubscriber(client mqtt.Client, hub *ws.Hub) {
	topic := "mesincuci/+/status"
	token := client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
		log.Printf("Terima status dari ESP32: %s", string(m.Payload()))

		// 1. Parse payload (Misal: {"mac": "AA:BB", "state": "WASHING", "rem_time": 10})
		var status map[string]interface{}
		if err := json.Unmarshal(m.Payload(), &status); err != nil {
			log.Printf("Gagal parse MQTT status: %v", err)
			return
		}

		// 2. Teruskan langsung ke semua Client WebSocket (Aplikasi Android)
		hub.BroadcastStatus(status)
	})

	token.Wait()
	log.Printf("Subscribed ke topik: %s", topic)
}
