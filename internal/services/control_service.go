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
	StopSession(macAddress string, sessionType string) error
}

type controlService struct {
	deviceRepo  repositories.DeviceRepository
	sessionRepo repositories.SessionRepository
	mqttClient  mqtt.Client
}

func NewControlService(dr repositories.DeviceRepository, sr repositories.SessionRepository, mc mqtt.Client) ControlService {
	return &controlService{deviceRepo: dr, sessionRepo: sr, mqttClient: mc}
}

type mqttPayload struct {
	Command  string `json:"cmd"`
	Duration int    `json:"dur"`
}

func (s *controlService) StartSession(macAddress string, sessionType string, duration int) (*models.WashSession, error) {
	_, err := s.deviceRepo.FindByMac(macAddress)
	if err != nil {
		return nil, errors.New("device tidak ditemukan")
	}

	cmd := "wash"
	if sessionType == "SPIN" {
		cmd = "spin"
	}

	payload, _ := json.Marshal(mqttPayload{Command: cmd, Duration: duration})
	topic := "mesincuci/" + macAddress + "/command"

	if token := s.mqttClient.Publish(topic, 1, false, payload); token.Wait() && token.Error() != nil {
		return nil, errors.New("gagal mengirim perintah ke device via MQTT")
	}

	session := &models.WashSession{
		DeviceID:        macAddress,
		SessionType:     sessionType,
		DurationMinutes: duration,
		StartTime:       time.Now(),
		Status:          "IN_PROGRESS",
	}

	s.sessionRepo.Create(session)

	if sessionType == "WASH" {
		s.deviceRepo.UpdateWashStatus(macAddress, "WASHING")
	} else if sessionType == "SPIN" {
		s.deviceRepo.UpdateSpinStatus(macAddress, "SPINNING")
	}

	return session, nil
}

func (s *controlService) StopSession(macAddress string, sessionType string) error {
	// 1. Tentukan command MQTT berdasarkan tipe
	cmd := "stop_wash"
	if sessionType == "SPIN" {
		cmd = "stop_spin"
	}

	payload, _ := json.Marshal(mqttPayload{Command: cmd, Duration: 0})
	topic := "mesincuci/" + macAddress + "/command"

	if token := s.mqttClient.Publish(topic, 1, false, payload); token.Wait() && token.Error() != nil {
		return errors.New("gagal mengirim perintah stop ke " + sessionType)
	}

	// 2. Update Database Sesi
	now := time.Now()
	s.sessionRepo.UpdateStatusByDeviceMAC(macAddress, "INTERRUPTED", &now)

	// 3. Update Status Device secara spesifik
	if sessionType == "WASH" {
		return s.deviceRepo.UpdateWashStatus(macAddress, "IDLE")
	} else {
		return s.deviceRepo.UpdateSpinStatus(macAddress, "IDLE")
	}
}
