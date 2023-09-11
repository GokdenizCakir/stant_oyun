package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GokdenizCakir/stant_oyun/src/dto"
	"github.com/GokdenizCakir/stant_oyun/src/models"
	"github.com/GokdenizCakir/stant_oyun/src/services"
	"github.com/GokdenizCakir/stant_oyun/src/utils"
	"github.com/gin-gonic/gin"
)

type PlayerController struct {
	playerService *services.PlayerService
}

type JWTPlayerData struct {
	ID        uint
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

	playerBody.Phone = strings.TrimPrefix(playerBody.Phone, "0")
	if !utils.IsPhoneNumber(playerBody.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz telefon numarası"})
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

	access_token, err := utils.GenerateJWT(JWTPlayerData{ID: player.ID, Questions: emptyQuestionsData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", access_token, 3600, "/", "https://skylabmilyoneryeni.vercel.app/", true, true)
	c.JSON(http.StatusCreated, gin.H{"jwt": access_token})
}

func (p *PlayerController) LogoutPlayer(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "https://skylabmilyoneryeni.vercel.app/", true, true)
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
