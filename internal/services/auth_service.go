package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeenarief/smart-washer-backend/internal/models"
	"github.com/zeenarief/smart-washer-backend/internal/repositories"
)

type AuthService interface {
	RegisterUser(username, password string) (*models.User, error)
	LoginUser(username, password string) (string, string, error)
	RefreshAccessToken(refreshToken string) (string, error)
}

type authService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) RegisterUser(username, password string) (*models.User, error) {
	// 1. Hash Password menggunakan Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	// 2. Simpan User
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("username mungkin sudah digunakan")
	}

	return user, nil
}

func (s *authService) LoginUser(username, password string) (string, string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", errors.New("username atau password salah")
	}

	accessToken, err := s.generateToken(user.ID, user.Username, "JWT_SECRET", "JWT_ACCESS_TOKEN_EXPIRE", 15)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateToken(user.ID, user.Username, "JWT_REFRESH_SECRET", "JWT_REFRESH_TOKEN_EXPIRE", 10080)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshAccessToken(refreshTokenStr string) (string, error) {
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	// 1. Validasi Refresh Token
	token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("refresh token tidak valid atau sudah kedaluwarsa")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("klaim token tidak valid")
	}

	userID := claims["user_id"].(string)
	username := claims["username"].(string)

	// 2. Buat Access Token baru
	newAccessToken, err := s.generateToken(userID, username, "JWT_SECRET", "JWT_ACCESS_TOKEN_EXPIRE", 15)
	if err != nil {
		return "", errors.New("gagal membuat access token baru")
	}

	return newAccessToken, nil
}

// Fungsi helper internal untuk generate token agar tidak duplikasi kode
func (s *authService) generateToken(userID, username, secretEnvKey, expireEnvKey string, defaultExpireMinutes int) (string, error) {
	secret := os.Getenv(secretEnvKey)
	expireStr := os.Getenv(expireEnvKey)

	expireMinutes, err := strconv.Atoi(expireStr)
	if err != nil {
		expireMinutes = defaultExpireMinutes
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
