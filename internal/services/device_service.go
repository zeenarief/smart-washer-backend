package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zeenarief/smart-washer-backend/internal/models"
	"github.com/zeenarief/smart-washer-backend/internal/repositories"
)

type DeviceService interface {
	RegisterDevice(userID, macAddress, name string) (*models.Device, error)
	GetStatus(macAddress string) (*models.Device, error)
	GetDevicesByUserID(userID string) ([]models.Device, error)
	UpdateDeviceName(userID string, macAddress string, newName string) error
	DeleteDevice(userID string, macAddress string) error
}

type deviceService struct {
	repo repositories.DeviceRepository
}

func (s *deviceService) DeleteDevice(userID string, macAddress string) error {
	return s.repo.Delete(userID, macAddress)
}

func (s *deviceService) UpdateDeviceName(userID string, macAddress string, newName string) error {
	// Anda bisa menambahkan validasi bisnis di sini jika perlu (misal: nama tidak boleh kosong)
	if newName == "" {
		return errors.New("nama perangkat tidak boleh kosong")
	}
	return s.repo.UpdateName(userID, macAddress, newName)
}

func (s *deviceService) GetDevicesByUserID(userID string) ([]models.Device, error) {
	return s.repo.FindByUserID(userID)
}

func NewDeviceService(repo repositories.DeviceRepository) DeviceService {
	return &deviceService{repo}
}

func (s *deviceService) RegisterDevice(userID, macAddress, name string) (*models.Device, error) {
	device := &models.Device{
		ID:         uuid.New().String(),
		UserID:     userID,
		MacAddress: macAddress,
		Name:       name,
		Status:     "IDLE",
		LastSeen:   time.Now(),
	}

	err := s.repo.Create(device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (s *deviceService) GetStatus(macAddress string) (*models.Device, error) {
	return s.repo.FindByMac(macAddress)
}
