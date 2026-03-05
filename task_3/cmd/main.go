package main

import (
	"log"
	"tasks/task_3/config"
	"tasks/task_3/routes"
)

func main() {
	// 1. Config yuklash
	cfg := config.Load()

	// 2. Database ulanish
	db := config.ConnectDB(cfg)

	// 3. Auto Migration
	config.Migrate(db)

	// 4. Router sozlash
	r := routes.SetupRouter(db, cfg)

	// 5. Server ishga tushirish
	log.Printf("🚀 Server %s portda ishlamoqda...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Server xatolik:", err)
	}
}
