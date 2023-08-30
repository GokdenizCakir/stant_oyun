package controllers

import (
	"net/http"

	"github.com/GokdenizCakir/stant_oyun/src/dto"
	"github.com/GokdenizCakir/stant_oyun/src/models"
	"github.com/GokdenizCakir/stant_oyun/src/services"
	"github.com/gin-gonic/gin"
)

type PlayerController struct {
	playerService *services.PlayerService
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

	c.JSON(http.StatusCreated, gin.H{"ID": player.ID})
}

func (p *PlayerController) GetScoreboard(c *gin.Context) {
	scoreboard, err := p.playerService.GetScoreboard()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": scoreboard})
}
