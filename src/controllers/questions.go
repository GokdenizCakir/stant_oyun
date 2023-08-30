package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/GokdenizCakir/stant_oyun/src/dto"
	"github.com/GokdenizCakir/stant_oyun/src/models"
	"github.com/GokdenizCakir/stant_oyun/src/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type QuestionController struct {
	QuestionService *services.QuestionService
}

func NewQuestionController(questionService *services.QuestionService) *QuestionController {
	return &QuestionController{
		QuestionService: questionService,
	}
}

func (q *QuestionController) CreateQuestion(c *gin.Context) {
	var questionBody *dto.CreateQuestionDto

	if err := c.ShouldBindJSON(&questionBody); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}

	var newQuestion = &models.Question{
		Value:      questionBody.Value,
		A:          questionBody.A,
		B:          questionBody.B,
		C:          questionBody.C,
		D:          questionBody.D,
		Answer:     questionBody.Answer,
		Difficulty: questionBody.Difficulty,
	}

	question, err := q.QuestionService.CreateQuestion(newQuestion)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"data": question})
}

func (q *QuestionController) GetQuestion(c *gin.Context) {

	difficulty := c.Param("difficulty")

	question, err := q.QuestionService.GetQuestion(difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	question.Answer = ""

	c.JSON(http.StatusOK, gin.H{"data": question})
}

func (q *QuestionController) AnswerQuestion(c *gin.Context) {
	var answerBody *dto.AnswerQuestionDto

	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := c.ShouldBindJSON(&answerBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	question, err := q.QuestionService.GetQuestionByID(uint(questionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if question.Answer == answerBody.Answer {
		score, err := q.QuestionService.IncreasePoints(answerBody.PlayerID, 1)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"data": true, "score": score})
	} else {
		score, err := q.QuestionService.IncreasePoints(answerBody.PlayerID, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"data": false, "score": score})
	}

}
