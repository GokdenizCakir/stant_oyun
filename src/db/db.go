package db

import (
	"fmt"

	"github.com/GokdenizCakir/stant_oyun/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	dsn := "host=trumpet.db.elephantsql.com user=fxwnvahh password=3AOt6mw6wPENHWSwHNcMyVDiVBhuGQyb dbname=fxwnvahh port=5432 sslmode=require TimeZone=Europe/London"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	err = DB.AutoMigrate(
		&models.Question{},
		&models.Player{},
		&models.JWT{},
	)
	if err != nil {
		panic("error migrating db")
	}

}
