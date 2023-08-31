package dto

type CreatePlayerDto struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"" binding:"required"`
}
