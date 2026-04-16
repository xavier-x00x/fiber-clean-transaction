package config

import (
	"fiber-clean-transaction/internal/domain/entity"
	"fmt"
	"log"
	"time"

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

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

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
		&entity.NuxtMenu{},
		&entity.Product{},
	)
	log.Println("✅ Database migrated")
}
