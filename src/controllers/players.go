package controllers

import (
	"net/http"

	"github.com/GokdenizCakir/stant_oyun/src/dto"
	"github.com/GokdenizCakir/stant_oyun/src/models"
	"github.com/GokdenizCakir/stant_oyun/src/services"
	"github.com/GokdenizCakir/stant_oyun/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlayerController struct {
	playerService *services.PlayerService
}

type JWTPlayerData struct {
	UUID      uuid.UUID
	Questions [][]int
}

func NewPlayerController(playerService *services.PlayerService) *PlayerController {
	return &PlayerController{playerService: playerService}
}

func (p *PlayerController) CreatePlayer(c *gin.Context) {
	var playerBody *dto.CreatePlayerDto

	if err := c.ShouldBindJSON(&playerBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var newPlayer = &models.Player{
		FullName: playerBody.FullName,
		Phone:    playerBody.Phone,
	}

	player, err := p.playerService.CreatePlayer(newPlayer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	access_token, err := utils.GenerateJWT(JWTPlayerData{UUID: player.ID, Questions: [][]int{{-1, -1}, {-1, -1}, {-1, -1}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.SetCookie("jwt", access_token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{"jwt": access_token})
}

func (p *PlayerController) GetScoreboard(c *gin.Context) {
	scoreboard, err := p.playerService.GetScoreboard()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": scoreboard})
}
