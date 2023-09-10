package services

import (
	"github.com/GokdenizCakir/stant_oyun/src/models"
	"gorm.io/gorm"
)

type QuestionService struct {
	DB       *gorm.DB
	Question *models.Question
}

func NewQuestionService(db *gorm.DB, question *models.Question) *QuestionService {
	return &QuestionService{
		DB:       db,
		Question: question,
	}
}

func (q *QuestionService) CreateQuestion(question *models.Question) (*models.Question, error) {
	if err := q.DB.Create(question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (q *QuestionService) GetQuestion(difficulty string) (*models.Question, error) {
	var question models.Question

	if err := q.DB.Where("difficulty = ?", difficulty).Order("times_asked asc").First(&question).Error; err != nil {
		return nil, err
	}

	question.TimesAsked += 1

	if err := q.DB.Save(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil

}

func (q *QuestionService) GetQuestionByID(id uint) (*models.Question, error) {
	var question models.Question

	if err := q.DB.Where("id = ?", id).First(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (p *QuestionService) IncreasePoints(ID float64, amount int) (int, error) {
	var player models.Player

	if err := p.DB.Where("id = ?", ID).First(&player).Error; err != nil {
		return 0, err
	}

	player.Score += amount

	if err := p.DB.Save(&player).Error; err != nil {
		return 0, err
	}

	return player.Score, nil
}
