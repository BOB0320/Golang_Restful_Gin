package initializers

import (
	"fmt"
	"log"

	"github.com/johnstewart0820/jurassic_park/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	DB.AutoMigrate(&models.Cage{}, &models.Dinosaur{})
	fmt.Println("🚀 Connected Successfully to the Database")
}
