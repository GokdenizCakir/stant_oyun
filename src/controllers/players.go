package controllers

import (
	"net/http"
	"os"
	"strconv"

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
		return
	}

	var newPlayer = &models.Player{
		FullName: playerBody.FullName,
		Phone:    playerBody.Phone,
	}

	player, err := p.playerService.CreatePlayer(newPlayer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questionCount, err := strconv.Atoi(os.Getenv("QUESTION_COUNT"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	emptyQuestionsData := make([][]int, questionCount)
	for i := range emptyQuestionsData {
		emptyQuestionsData[i] = []int{-1, -1}
	}

	access_token, err := utils.GenerateJWT(JWTPlayerData{UUID: player.ID, Questions: emptyQuestionsData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("jwt", access_token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{"jwt": access_token})
}

func (p *PlayerController) LogoutPlayer(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (p *PlayerController) GetScoreboard(c *gin.Context) {
	scoreboard, err := p.playerService.GetScoreboard()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scoreboard})
}
