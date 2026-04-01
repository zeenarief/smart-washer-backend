package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeenarief/smart-washer-backend/internal/services"
	"github.com/zeenarief/smart-washer-backend/pkg/response"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := response.Error("Format request tidak valid")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	user, err := h.service.RegisterUser(req.Username, req.Password)
	if err != nil {
		res := response.Error(err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Success("Registrasi berhasil", user)
	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest // Struct dari kode sebelumnya
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Format request tidak valid"))
		return
	}

	accessToken, refreshToken, err := h.service.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
		return
	}

	data := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	c.JSON(http.StatusOK, response.Success("Login berhasil", data))
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Refresh token harus disertakan"))
		return
	}

	newAccessToken, err := h.service.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		// HTTP 401 berarti di sisi Flutter harus otomatis diarahkan ke halaman Login
		c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
		return
	}

	data := map[string]string{
		"access_token": newAccessToken,
	}

	c.JSON(http.StatusOK, response.Success("Access token diperbarui", data))
}
