package models

import "time"

type Player struct {
	ID          uint      `gorm:"primaryKey"`
	FullName    string    `json:"full_name" gorm:"size:40"`
	Phone       string    `json:"phone" gorm:"unique"`
	IP          string    `json:"ip"` /* gorm:"unique" */
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Score       int       `json:"score" gorm:"default:0"`
	HasFinished bool      `json:"has_finished" gorm:"default:false"`
}
