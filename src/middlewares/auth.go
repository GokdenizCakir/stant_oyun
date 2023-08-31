package middlewares

import (
	"net/http"

	"github.com/GokdenizCakir/stant_oyun/src/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* token := c.GetHeader("Authorization") */
		token, err := c.Cookie("jwt")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization token",
			})
			return
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization token",
			})
			return
		}

		/* token = strings.Split(token, " ")[1] */

		data, err := utils.HandleJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set("user", data)
		c.Next()
	}
}
