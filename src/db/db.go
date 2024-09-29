package db

import (
	"fmt"

	"github.com/GokdenizCakir/stant_oyun/src/models"
	"github.com/GokdenizCakir/stant_oyun/src/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(config utils.Config) {
	var err error

	fmt.Println(config.DBSource)

	DB, err = gorm.Open(postgres.Open(config.DBSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	err = DB.AutoMigrate(
		&models.Question{},
		&models.Player{},
	)
	if err != nil {
		panic("error migrating db")
	}

}
