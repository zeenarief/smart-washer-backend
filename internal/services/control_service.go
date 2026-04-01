package services

import (
	"encoding/json"
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeenarief/smart-washer-backend/internal/models"
	"github.com/zeenarief/smart-washer-backend/internal/repositories"
)

type ControlService interface {
	StartSession(macAddress string, sessionType string, duration int) (*models.WashSession, error)
	StopSession(macAddress string) error
}

type controlService struct {
	deviceRepo  repositories.DeviceRepository
	sessionRepo repositories.SessionRepository
	mqttClient  mqtt.Client
}

func NewControlService(dr repositories.DeviceRepository, sr repositories.SessionRepository, mc mqtt.Client) ControlService {
	return &controlService{deviceRepo: dr, sessionRepo: sr, mqttClient: mc}
}

// Struct untuk payload MQTT internal
type mqttPayload struct {
	Command  string `json:"cmd"`
	Duration int    `json:"dur"`
}

func (s *controlService) StartSession(macAddress string, sessionType string, duration int) (*models.WashSession, error) {
	// 1. Cek apakah device ada
	_, err := s.deviceRepo.FindByMac(macAddress)
	if err != nil {
		return nil, errors.New("device tidak ditemukan")
	}

	// 2. Kirim perintah via MQTT
	cmd := "wash"
	if sessionType == "SPIN" {
		cmd = "spin"
	}

	payload, _ := json.Marshal(mqttPayload{Command: cmd, Duration: duration})
	topic := "mesincuci/" + macAddress + "/command"

	if token := s.mqttClient.Publish(topic, 1, false, payload); token.Wait() && token.Error() != nil {
		return nil, errors.New("gagal mengirim perintah ke device via MQTT")
	}

	// 3. Catat di DB
	session := &models.WashSession{
		DeviceID:        macAddress,
		SessionType:     sessionType,
		DurationMinutes: duration,
		StartTime:       time.Now(),
		Status:          "IN_PROGRESS",
	}

	s.sessionRepo.Create(session)

	// 4. Update status Device
	newStatus := "WASHING"
	if sessionType == "SPIN" {
		newStatus = "SPINNING"
	}
	s.deviceRepo.UpdateStatus(macAddress, newStatus)

	return session, nil
}

func (s *controlService) StopSession(macAddress string) error {
	// Kirim perintah stop via MQTT
	payload, _ := json.Marshal(mqttPayload{Command: "stop", Duration: 0})
	topic := "mesincuci/" + macAddress + "/command"

	if token := s.mqttClient.Publish(topic, 1, false, payload); token.Wait() && token.Error() != nil {
		return errors.New("gagal mengirim perintah stop via MQTT")
	}

	// Update DB: Hentikan sesi yang aktif dan set device ke IDLE
	now := time.Now()
	s.sessionRepo.UpdateStatusByDeviceMAC(macAddress, "INTERRUPTED", &now)
	s.deviceRepo.UpdateStatus(macAddress, "IDLE")

	return nil
}
