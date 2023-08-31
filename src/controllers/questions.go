package controllers

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/GokdenizCakir/stant_oyun/src/dto"
	"github.com/GokdenizCakir/stant_oyun/src/models"
	"github.com/GokdenizCakir/stant_oyun/src/services"
	"github.com/GokdenizCakir/stant_oyun/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	JWTData := c.MustGet("user")
	JWTQuestions := JWTData.(map[string]interface{})["Questions"].([]interface{})

	var questionIndex int
	var difficulty string

	for i := range JWTQuestions {
		if JWTQuestions[i].([]interface{})[1].(float64) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You already answered wrong"})
			return
		} else if JWTQuestions[i].([]interface{})[0].(float64) != -1 && JWTQuestions[i].([]interface{})[1].(float64) == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You haven't answered yet"})
			return
		} else if JWTQuestions[i].([]interface{})[0].(float64) == -1 && JWTQuestions[i].([]interface{})[1].(float64) == -1 {
			questionIndex = i
			break
		}
	}

	difficulty = fmt.Sprintf("%.0f", math.Ceil(float64(questionIndex+1)/2))
	question, err := q.QuestionService.GetQuestion(difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question.Answer = ""
	JWTQuestions[questionIndex] = []interface{}{float64(question.ID), -1}

	err = utils.UpdateJWT(c, "Questions", JWTQuestions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": question})
}

func (q *QuestionController) AnswerQuestion(c *gin.Context) {
	var answerBody *dto.AnswerQuestionDto

	JWTData := c.MustGet("user")
	JWTPlayerID, err := uuid.Parse(JWTData.(map[string]interface{})["UUID"].(string))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	JWTQuestions := JWTData.(map[string]interface{})["Questions"].([]interface{})
	var questionIndex int
	var questionID int

	for i := range JWTQuestions {
		if JWTQuestions[i].([]interface{})[1].(float64) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You already answered wrong"})
			return
		} else if JWTQuestions[i].([]interface{})[1].(float64) == -1 {
			questionIndex = i
			break
		}
	}

	questionID = int(JWTQuestions[questionIndex].([]interface{})[0].(float64))
	if questionID == -1 {
		fmt.Println("You haven't seen the question yet")
		c.JSON(http.StatusBadRequest, gin.H{"error": "You haven't seen the question yet"})
		return
	}

	if err := c.ShouldBindJSON(&answerBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := q.QuestionService.GetQuestionByID(uint(questionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if question.Answer == answerBody.Answer {
		score, err := q.QuestionService.IncreasePoints(JWTPlayerID, 1)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		JWTQuestions[questionIndex] = []interface{}{questionID, 1}
		utils.UpdateJWT(c, "Questions", JWTQuestions)

		c.JSON(http.StatusOK, gin.H{"data": true, "score": score})
		return
	} else {
		score, err := q.QuestionService.IncreasePoints(JWTPlayerID, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		JWTQuestions[questionIndex] = []interface{}{questionID, 0}
		utils.UpdateJWT(c, "Questions", JWTQuestions)

		c.JSON(http.StatusOK, gin.H{"data": false, "score": score})
	}

}
