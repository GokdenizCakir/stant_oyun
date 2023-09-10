package models

type Player struct {
	ID          uint   `gorm:"primaryKey"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone" gorm:"unique"`
	Score       int    `json:"score" gorm:"default:0"`
	HasFinished bool   `json:"has_finished" gorm:"default:false"`
}
