package routes

import (
	"github.com/GokdenizCakir/stant_oyun/src/controllers"
	"github.com/GokdenizCakir/stant_oyun/src/db"
	"github.com/GokdenizCakir/stant_oyun/src/models"
	"github.com/GokdenizCakir/stant_oyun/src/services"
	"github.com/gin-gonic/gin"
)

func QuestionRouter(r *gin.RouterGroup) {
	questionRouter := r.Group("/questions")
	{
		questionService := services.NewQuestionService(db.DB, &models.Question{})
		questionController := controllers.NewQuestionController(questionService)

		questionRouter.
			POST("", questionController.CreateQuestion).
			GET(":difficulty", questionController.GetQuestion).
			POST(":id", questionController.AnswerQuestion)
	}
}
