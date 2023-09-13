package dto

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
	Answer    string `json:"answer" binding:"required"`
	HasGaveUp bool   `json:"hasGaveUp" binding:"required"`
}
