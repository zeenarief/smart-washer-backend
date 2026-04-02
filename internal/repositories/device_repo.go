package repositories

import (
	"errors"

	"github.com/zeenarief/smart-washer-backend/internal/models"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	Create(device *models.Device) error
	FindByMac(mac string) (*models.Device, error)
	UpdateWashStatus(mac string, status string) error
	UpdateSpinStatus(mac string, status string) error
	UpdateAllStatus(mac string, washStatus string, spinStatus string) error
	FindByUserID(userID string) ([]models.Device, error)
	UpdateName(userID string, macAddress string, newName string) error
	Delete(userID string, macAddress string) error
}

type deviceRepository struct {
	db *gorm.DB
}

// Implementasi fungsi-fungsi baru
func (r *deviceRepository) UpdateWashStatus(mac string, status string) error {
	return r.db.Model(&models.Device{}).Where("mac_address = ?", mac).Update("wash_status", status).Error
}

func (r *deviceRepository) UpdateSpinStatus(mac string, status string) error {
	return r.db.Model(&models.Device{}).Where("mac_address = ?", mac).Update("spin_status", status).Error
}

func (r *deviceRepository) UpdateAllStatus(mac string, washStatus string, spinStatus string) error {
	return r.db.Model(&models.Device{}).Where("mac_address = ?", mac).Updates(map[string]interface{}{
		"wash_status": washStatus,
		"spin_status": spinStatus,
	}).Error
}

// Implementasi fungsi lama (tetap sama)
func (r *deviceRepository) Delete(userID string, macAddress string) error {
	result := r.db.Unscoped().Where("user_id = ? AND mac_address = ?", userID, macAddress).Delete(&models.Device{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("perangkat tidak ditemukan atau Anda tidak memiliki akses")
	}
	return nil
}

func (r *deviceRepository) UpdateName(userID string, macAddress string, newName string) error {
	result := r.db.Model(&models.Device{}).Where("user_id = ? AND mac_address = ?", userID, macAddress).Update("name", newName)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("perangkat tidak ditemukan atau Anda tidak memiliki akses")
	}
	return nil
}

func (r *deviceRepository) FindByUserID(userID string) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.Where("user_id = ?", userID).Find(&devices).Error
	return devices, err
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db}
}

func (r *deviceRepository) Create(device *models.Device) error {
	return r.db.Create(device).Error
}

func (r *deviceRepository) FindByMac(mac string) (*models.Device, error) {
	var device models.Device
	err := r.db.Where("mac_address = ?", mac).First(&device).Error
	return &device, err
}
