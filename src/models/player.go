package models

import "github.com/google/uuid"

type Player struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	FullName string    `json:"full_name" gorm:"unique"`
	Phone    string    `json:"phone" gorm:"unique"`
	Score    int       `json:"score" gorm:"default:0"`
}