package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeenarief/smart-washer-backend/internal/handlers"
	"github.com/zeenarief/smart-washer-backend/internal/middleware"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, deviceHandler *handlers.DeviceHandler, controlHandler *handlers.ControlHandler) {

	v1 := router.Group("/api/v1")
	{
		// Public Routes (Tidak perlu token)
		v1.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "Pong!"}) })

		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected Routes (Wajib memiliki token JWT yang valid)
		protected := v1.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			deviceGroup := protected.Group("/device")
			{
				deviceGroup.GET("/", deviceHandler.GetUserDevices)
				deviceGroup.POST("/register", deviceHandler.RegisterDevice)
				deviceGroup.GET("/status/:mac_address", deviceHandler.GetDeviceStatus)
				deviceGroup.PUT("/:mac_address", deviceHandler.UpdateDevice)
				deviceGroup.DELETE("/:mac_address", deviceHandler.DeleteDevice)
			}

			controlGroup := protected.Group("/control")
			{
				controlGroup.POST("/wash", controlHandler.StartWash)
				controlGroup.POST("/spin", controlHandler.StartSpin)
				controlGroup.POST("/stop/wash/:mac_address", controlHandler.StopWash)
				controlGroup.POST("/stop/spin/:mac_address", controlHandler.StopSpin)
			}
		}
	}
}
