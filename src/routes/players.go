package routes

import (
	"github.com/GokdenizCakir/stant_oyun/src/controllers"
	"github.com/GokdenizCakir/stant_oyun/src/db"
	"github.com/GokdenizCakir/stant_oyun/src/middlewares"
	"github.com/GokdenizCakir/stant_oyun/src/models"
	"github.com/GokdenizCakir/stant_oyun/src/services"
	"github.com/gin-gonic/gin"
)

func PlayerRouter(r *gin.RouterGroup) {
	playerRouter := r.Group("/players")
	{
		playerService := services.NewPlayerService(db.DB, &models.Player{})
		playerController := controllers.NewPlayerController(playerService)

		playerRouter.
			POST("", playerController.CreatePlayer).
			GET("", playerController.GetScoreboard).
			POST("/logout", middlewares.AuthMiddleware(), playerController.LogoutPlayer)
	}
}
