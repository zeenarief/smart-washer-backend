package models

import (
	"time"
)

// Device merepresentasikan mesin cuci ESP32
type Device struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID     string    `gorm:"type:varchar(36);index;not null" json:"user_id"` // KUNCI KEPEMILIKAN
	MacAddress string    `gorm:"uniqueIndex;type:varchar(50);not null" json:"mac_address"`
	Name       string    `gorm:"type:varchar(100)" json:"name"`
	WashStatus string    `gorm:"type:varchar(20);default:'IDLE'" json:"wash_status"` // IDLE, WASHING
	SpinStatus string    `gorm:"type:varchar(20);default:'IDLE'" json:"spin_status"` // IDLE, SPINNING
	LastSeen   time.Time `json:"last_seen"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// WashSession mencatat riwayat penggunaan mesin cuci (Tidak diubah)
type WashSession struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID        string     `gorm:"type:varchar(36);index" json:"device_id"`
	SessionType     string     `gorm:"type:varchar(20)" json:"session_type"` // WASH, SPIN
	DurationMinutes int        `json:"duration_minutes"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	Status          string     `gorm:"type:varchar(20);default:'IN_PROGRESS'" json:"status"`
}
