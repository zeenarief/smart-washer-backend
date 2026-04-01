package repositories

import (
	"github.com/zeenarief/smart-washer-backend/internal/models"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(session *models.WashSession) error
	UpdateStatusByDeviceMAC(mac string, status string, endTime interface{}) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db}
}

func (r *sessionRepository) Create(session *models.WashSession) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) UpdateStatusByDeviceMAC(mac string, status string, endTime interface{}) error {
	return r.db.Model(&models.WashSession{}).
		Where("device_id = ? AND status = ?", mac, "IN_PROGRESS").
		Updates(map[string]interface{}{
			"status":   status,
			"end_time": endTime,
		}).Error
}
