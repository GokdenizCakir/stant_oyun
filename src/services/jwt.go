package services

import (
	"github.com/GokdenizCakir/stant_oyun/src/db"
	"github.com/GokdenizCakir/stant_oyun/src/models"
)

func InsertJWT(jwt *models.JWT) error {
	if err := db.DB.Create(*jwt).Error; err != nil {
		return err
	}

	return nil
}
