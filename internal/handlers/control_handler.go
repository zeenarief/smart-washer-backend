package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeenarief/smart-washer-backend/internal/services"
)

type ControlHandler struct {
	service services.ControlService
}

func NewControlHandler(service services.ControlService) *ControlHandler {
	return &ControlHandler{service}
}

type ControlRequest struct {
	MacAddress string `json:"mac_address" binding:"required"`
	Duration   int    `json:"duration" binding:"required"`
}

func (h *ControlHandler) StartWash(c *gin.Context) {
	h.processSession(c, "WASH")
}

func (h *ControlHandler) StartSpin(c *gin.Context) {
	h.processSession(c, "SPIN")
}

// Fungsi helper internal untuk mengurangi duplikasi kode
func (h *ControlHandler) processSession(c *gin.Context, sessionType string) {
	var req ControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	session, err := h.service.StartSession(req.MacAddress, sessionType, req.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Perintah " + sessionType + " berhasil dikirim",
		"data":    session,
	})
}

func (h *ControlHandler) StopMachine(c *gin.Context) {
	macAddress := c.Param("mac_address")

	err := h.service.StopSession(macAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesin cuci berhasil dihentikan"})
}
