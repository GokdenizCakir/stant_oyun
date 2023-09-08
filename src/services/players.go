package services

import (
	"github.com/GokdenizCakir/stant_oyun/src/models"

	"gorm.io/gorm"
)

type PlayerService struct {
	DB     *gorm.DB
	Player *models.Player
}

type PlayerScore struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
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

func (p *PlayerService) GetScoreboard() ([]PlayerScore, error) {
	var players []models.Player
	var playerScores []PlayerScore

	if err := p.DB.Order("score desc").Limit(10).Find(&players).Error; err != nil {
		return nil, err
	}

	for _, player := range players {
		playerScores = append(playerScores, PlayerScore{
			Name:  player.FullName,
			Score: player.Score,
		})
	}

	return playerScores, nil
}
