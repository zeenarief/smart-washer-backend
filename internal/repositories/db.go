package repositories

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/zeenarief/smart-washer-backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// Mengambil konfigurasi dari Environment Variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke MySQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan instance sql.DB: %v", err)
	}

	// Parse nilai dari .env dengan fallback default jika kosong
	maxOpenConns, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_IDLE_CONNS", "5"))
	maxLifetimeSeconds, _ := strconv.Atoi(getEnvOrDefault("DB_CONN_MAX_LIFETIME", "300"))

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetimeSeconds) * time.Second)

	log.Println("Koneksi MySQL & Connection Pool berhasil disiapkan!")

	// Auto-Migrate: Membuat/mengupdate tabel otomatis
	err = db.AutoMigrate(&models.Device{}, &models.WashSession{}, &models.User{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	return db
}

// Helper kecil untuk fallback environment variable
func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
