package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zeenarief/smart-washer-backend/pkg/response"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			res := response.Error("Akses ditolak. Token tidak ditemukan.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		// 2. Format header harus "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			res := response.Error("Format token tidak valid. Gunakan: Bearer <token>")
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		tokenString := parts[1]
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "rahasia_default_untuk_dev"
		}

		// 3. Validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signin tidak valid")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			res := response.Error("Token tidak valid atau sudah kedaluwarsa")
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		// 4. Ekstrak payload (claims) dan simpan di context Gin
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
		}

		// 5. Lanjutkan ke Handler berikutnya
		c.Next()
	}
}
