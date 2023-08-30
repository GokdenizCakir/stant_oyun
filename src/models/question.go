package models

type Question struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Value      string `json:"value" gorm:"not null"`
	A          string `json:"a" gorm:"not null"`
	B          string `json:"b" gorm:"not null"`
	C          string `json:"c" gorm:"not null"`
	D          string `json:"d" gorm:"not null"`
	Answer     string `json:"answer" gorm:"not null"`
	TimesAsked int    `json:"times_asked" gorm:"not null default:0"`
	Difficulty int    `json:"difficulty" gorm:"not null"`
}
