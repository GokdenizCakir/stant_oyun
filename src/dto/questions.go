package dto

import "github.com/google/uuid"

type CreateQuestionDto struct {
	Value      string `json:"value" binding:"required"`
	A          string `json:"a" binding:"required"`
	B          string `json:"b" binding:"required"`
	C          string `json:"c" binding:"required"`
	D          string `json:"d" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
	Difficulty int    `json:"difficulty" binding:"required"`
}

type AnswerQuestionDto struct {
	Answer     string    `json:"answer" binding:"required"`
	PlayerID   uuid.UUID `json:"player_id" binding:"required"`
}
