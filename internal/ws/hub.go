package ws

import (
	"encoding/json"
	"log"
)

// Hub mengelola set klien yang aktif dan menyebarkan pesan ke klien.
type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.Println("WebSocket Client Terhubung")
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Println("WebSocket Client Terputus")
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// BroadcastStatus helper untuk mengirim struct sebagai JSON ke semua client WS
func (h *Hub) BroadcastStatus(data interface{}) {
	payload, err := json.Marshal(data)
	if err == nil {
		h.Broadcast <- payload
	}
}
