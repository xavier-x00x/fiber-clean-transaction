package config

import (
	"fiber-clean-transaction/internal/domain/entity"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ConfigApp.DBUser,
		ConfigApp.DBPass,
		ConfigApp.DBHost,
		ConfigApp.DBPort,
		ConfigApp.DBName,
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = database
	log.Println("✅ Connected to database")

	DB.AutoMigrate(
		&entity.User{},
		&entity.Store{},
		&entity.Category{},
		&entity.Unit{},
		&entity.Tax{},
		&entity.Permission{},
		&entity.Role{},
		&entity.RolePermission{},
		&entity.NumberSequence{},
	)
	log.Println("✅ Database migrated")
}
