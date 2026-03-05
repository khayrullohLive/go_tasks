package config

import (
	"fmt"
	"log"
	"os"
	"tasks/task_3/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config - barcha environment o'zgaruvchilar
type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
}

// Load - .env fayldan config yuklaydi
func Load() *Config {
	// .env fayl mavjud bo'lsa yuklaydi (production'da kerak emas)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env fayl topilmadi, environment variables ishlatiladi")
	}

	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "macbook"),
		DBPass:    getEnv("DB_PASS", "postgres"),
		DBName:    getEnv("DB_NAME", "blogdb"),
		JWTSecret: getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
	}
}

// ConnectDB - PostgreSQL ga ulanadi
func ConnectDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tashkent",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // SQL loglarini ko'rsatadi
	})
	if err != nil {
		log.Fatal("❌ Database ulanishda xato:", err)
	}

	log.Println("✅ Database ga muvaffaqiyatli ulandi!")
	return db
}

// Migrate - barcha modellarni avtomatik migratsiya qiladi
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
		&models.Category{},
		&models.Like{},
		&models.Follow{},
	)
	if err != nil {
		log.Fatal("❌ Migration xatolik:", err)
	}
	log.Println("✅ Migration muvaffaqiyatli bajarildi!")
}

// getEnv - environment variable ni oladi, agar yo'q bo'lsa default qaytaradi
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
