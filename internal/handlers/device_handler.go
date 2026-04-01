package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeenarief/smart-washer-backend/internal/services"
	"github.com/zeenarief/smart-washer-backend/pkg/response"
)

type DeviceHandler struct {
	service services.DeviceService
}

type RegisterDeviceRequest struct {
	MacAddress string `json:"mac_address" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

type UpdateDeviceRequest struct {
	Name string `json:"name" binding:"required"`
}

func NewDeviceHandler(service services.DeviceService) *DeviceHandler {
	return &DeviceHandler{service}
}

func (h *DeviceHandler) RegisterDevice(c *gin.Context) {
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil user_id dari middleware JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error("User tidak terautentikasi"))
		return
	}

	device, err := h.service.RegisterDevice(userID.(string), req.MacAddress, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Gagal mendaftarkan device. MAC mungkin sudah dipakai."))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Device berhasil didaftarkan", "data": device})
}

func (h *DeviceHandler) GetDeviceStatus(c *gin.Context) {
	macAddress := c.Param("mac_address")

	device, err := h.service.GetStatus(macAddress)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": device})
}

func (h *DeviceHandler) GetUserDevices(c *gin.Context) {
	userID, _ := c.Get("user_id")
	devices, err := h.service.GetDevicesByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Gagal mengambil daftar perangkat"))
		return
	}
	c.JSON(http.StatusOK, response.Success("Daftar perangkat berhasil diambil", devices))
}

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	macAddress := c.Param("mac_address")
	var req UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Format request tidak valid"))
		return
	}

	userID, _ := c.Get("user_id")

	// Panggil service untuk update nama di DB (Pastikan validasi userID agar tidak bisa mengubah milik orang lain)
	err := h.service.UpdateDeviceName(userID.(string), macAddress, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Gagal memperbarui perangkat"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Perangkat berhasil diperbarui", nil))
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	macAddress := c.Param("mac_address")
	userID, _ := c.Get("user_id")

	// Panggil service untuk hapus dari DB (Pastikan validasi userID)
	err := h.service.DeleteDevice(userID.(string), macAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Gagal menghapus perangkat"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Perangkat berhasil dihapus", nil))
}
