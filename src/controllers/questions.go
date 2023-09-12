package controllers

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/GokdenizCakir/stant_oyun/src/dto"
	"github.com/GokdenizCakir/stant_oyun/src/models"
	"github.com/GokdenizCakir/stant_oyun/src/services"
	"github.com/GokdenizCakir/stant_oyun/src/utils"
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
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
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
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": question})
}

func (q *QuestionController) GetQuestion(c *gin.Context) {
	JWTData := c.MustGet("user")
	JWTQuestions := JWTData.(map[string]interface{})["Questions"].([]interface{})

	questionIndex := -1
	var difficulty string

	for i := range JWTQuestions {
		if JWTQuestions[i].([]interface{})[1].(float64) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You already answered wrong"})
			return
		} else if JWTQuestions[i].([]interface{})[0].(float64) != -1 && JWTQuestions[i].([]interface{})[1].(float64) == -1 {
			/* Viewed question but refreshed page without answering */
			question, err := q.QuestionService.GetQuestionByID(uint(JWTQuestions[i].([]interface{})[0].(float64)))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			question.Answer = ""
			lastViewedAt := time.Unix(int64(JWTData.(map[string]interface{})["LastViewedAt"].(float64)), 0)
			questionSeconds, err := strconv.ParseFloat(os.Getenv("QUESTION_SECONDS"), 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			timeLeft := questionSeconds - time.Since(lastViewedAt).Seconds()
			c.JSON(http.StatusOK, gin.H{"data": question, "timeLeft": timeLeft})
			return
		} else if JWTQuestions[i].([]interface{})[0].(float64) == -1 && JWTQuestions[i].([]interface{})[1].(float64) == -1 {
			questionIndex = i
			break
		}
	}

	if questionIndex == -1 {
		c.JSON(http.StatusOK, gin.H{"data": "You won"})
		return
	}

	difficulty = strconv.Itoa(questionIndex + 1)
	question, err := q.QuestionService.GetQuestion(difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question.Answer = ""
	JWTQuestions[questionIndex] = []interface{}{float64(question.ID), -1}

	err = utils.UpdateJWT(c, "Questions", JWTQuestions, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": question})
}

func (q *QuestionController) AnswerQuestion(c *gin.Context) {
	var answerBody *dto.AnswerQuestionDto

	if err := c.ShouldBindJSON(&answerBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	JWTData := c.MustGet("user")
	JWTPlayerID := JWTData.(map[string]interface{})["ID"].(float64)
	lastViewedAt := JWTData.(map[string]interface{})["LastViewedAt"].(float64)

	JWTQuestions := JWTData.(map[string]interface{})["Questions"].([]interface{})
	questionIndex := -1
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

	if questionIndex == -1 {
		c.JSON(http.StatusOK, gin.H{"data": "You already won"})
		return
	}

	questionID = int(JWTQuestions[questionIndex].([]interface{})[0].(float64))
	if questionID == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You haven't seen the question yet"})
		return
	}

	questionSeconds, err := strconv.Atoi(os.Getenv("QUESTION_SECONDS"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := q.QuestionService.GetQuestionByID(uint(questionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if time.Now().Unix() > (int64(lastViewedAt) + int64(questionSeconds)) {
		JWTQuestions[questionIndex] = []interface{}{questionID, 0}
		utils.UpdateJWT(c, "Questions", JWTQuestions, false)

		c.JSON(http.StatusOK, gin.H{"answer": question.Answer})
		return
	}

	playerScore, playerHasFinished, err := q.QuestionService.GetPlayerStatus(JWTPlayerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if (questionIndex < playerScore) || playerHasFinished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already answered this question."})
		return
	}

	if question.Answer == answerBody.Answer {
		playerScoreAfterAnswer, err := q.QuestionService.IncreasePoints(JWTPlayerID, 1)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		questionCount, err := strconv.Atoi(os.Getenv("QUESTION_COUNT"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		JWTQuestions[questionIndex] = []interface{}{questionID, 1}
		utils.UpdateJWT(c, "Questions", JWTQuestions, false)

		if playerScoreAfterAnswer == questionCount {
			c.JSON(http.StatusOK, gin.H{"answer": question.Answer, "winner": true})
			return
		}
		c.JSON(http.StatusOK, gin.H{"answer": question.Answer})
		return
	} else {
		_, err := q.QuestionService.IncreasePoints(JWTPlayerID, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		JWTQuestions[questionIndex] = []interface{}{questionID, 0}
		utils.UpdateJWT(c, "Questions", JWTQuestions, false)

		c.JSON(http.StatusOK, gin.H{"answer": question.Answer})
	}

}
