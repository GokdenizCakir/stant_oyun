package services

import (
	"github.com/GokdenizCakir/stant_oyun/src/models"

	"gorm.io/gorm"
)

type PlayerService struct {
	DB     *gorm.DB
	Player *models.Player
}

func NewPlayerService(db *gorm.DB, player *models.Player) *PlayerService {
	return &PlayerService{
		DB:     db,
		Player: player,
	}
}

func (p *PlayerService) CreatePlayer(player *models.Player) (*models.Player, error) {
	if err := p.DB.Create(player).Error; err != nil {
		return nil, err
	}

	return player, nil
}

func (p *PlayerService) GetScoreboard() ([]models.Player, error) {
	var players []models.Player

	if err := p.DB.Order("score desc").Find(&players).Error; err != nil {
		return nil, err
	}

	return players, nil
}
