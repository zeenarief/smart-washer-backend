package response

// WebResponse adalah struktur standar untuk semua balasan API
type WebResponse struct {
	Status  bool        `json:"status"`         // true untuk sukses, false untuk error
	Message string      `json:"message"`        // Pesan deskriptif
	Data    interface{} `json:"data,omitempty"` // Data payload (dihilangkan dari JSON jika nil)
}

// Success adalah helper untuk membuat response sukses dengan cepat
func Success(message string, data interface{}) WebResponse {
	return WebResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
}

// Error adalah helper untuk membuat response error dengan cepat
func Error(message string) WebResponse {
	return WebResponse{
		Status:  false,
		Message: message,
		Data:    nil,
	}
}
