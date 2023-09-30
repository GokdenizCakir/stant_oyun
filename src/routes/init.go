package routes

import (
	"github.com/GokdenizCakir/stant_oyun/src/middlewares"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func Init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.Default()
	Router.Use(middlewares.CORSMiddleware())

	api := Router.Group("api")
	{
		v1 := api.Group("v1")
		{
			QuestionRouter(v1)
			PlayerRouter(v1)
		}
	}

}
