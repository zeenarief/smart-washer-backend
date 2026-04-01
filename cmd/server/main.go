package main

import (
	"log"
	_ "net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/zeenarief/smart-washer-backend/internal/handlers"
	"github.com/zeenarief/smart-washer-backend/internal/mqttclient"
	"github.com/zeenarief/smart-washer-backend/internal/repositories"
	"github.com/zeenarief/smart-washer-backend/internal/routes"
	"github.com/zeenarief/smart-washer-backend/internal/services"
	"github.com/zeenarief/smart-washer-backend/internal/ws"
)

func main() {
	godotenv.Load()

	// 1. Init Database & MQTT
	db := repositories.InitDB()
	mqttClient := mqttclient.InitMQTT()

	// 2. Init WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// 3. Setup MQTT Subscriber (Mendengarkan ESP32 -> Kirim ke WS)
	mqttclient.SetupSubscriber(mqttClient, hub)

	// 4. Dependency Injection
	userRepo := repositories.NewUserRepository(db)
	deviceRepo := repositories.NewDeviceRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)

	authService := services.NewAuthService(userRepo)
	deviceService := services.NewDeviceService(deviceRepo)
	controlService := services.NewControlService(deviceRepo, sessionRepo, mqttClient)

	authHandler := handlers.NewAuthHandler(authService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)
	controlHandler := handlers.NewControlHandler(controlService)

	// 5. Gin Router
	router := gin.Default()

	// Endpoint Khusus WebSocket untuk Flutter
	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c.Writer, c.Request)
	})

	routes.SetupRoutes(router, authHandler, deviceHandler, controlHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Backend Smart Washer siap di port %s", port)
	router.Run(":" + port)
}
