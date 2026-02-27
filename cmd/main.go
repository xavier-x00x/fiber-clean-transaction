package main

import (
	"fiber-clean-transaction/internal/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()      // baca .env ke dalam config.Config
	config.ConnectDatabase() // koneksi DB pakai GORM dan simpan ke config.DB

	app := fiber.New()
	db := config.DB

	config.Bootstrap(&config.BootstrapConfig{
		DB:  db,
		App: app,
	})

	port := config.ConfigApp.AppPort
	app.Listen(":" + port)
}

// go mod init github.com/username/fiber-clean-transaction
// go mod tidy
// go run cmd/main.go
