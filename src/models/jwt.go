package models

type JWT struct {
	JWT string `gorm:"unique;not null"`
}
