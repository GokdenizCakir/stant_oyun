package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	HasGaveUp bool
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

	if len(playerBody.FullName) < 3 || len(playerBody.FullName) > 40 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "İsim 3-40 karakter arasında olmalıdır"})
		return
	}

	playerBody.Phone = strings.TrimPrefix(playerBody.Phone, "0")
	if !utils.IsPhoneNumber(playerBody.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz telefon numarası"})
		return
	}

	loc, _ := time.LoadLocation("Europe/Istanbul")

	var newPlayer = &models.Player{
		FullName:     playerBody.FullName,
		Phone:        playerBody.Phone,
		IP:           c.ClientIP(),
		CreationTime: time.Now().In(loc).Format("2006-01-02 15:04:05"),
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

	access_token, err := utils.GenerateJWT(JWTPlayerData{ID: player.ID, Questions: emptyQuestionsData, HasGaveUp: false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", access_token, 3600, "/", "egehan.dev", true, true)
	c.JSON(http.StatusCreated, gin.H{"jwt": access_token})
}

func (p *PlayerController) LogoutPlayer(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "egehan.dev", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (p *PlayerController) GetScoreboard(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scoreboard, err := p.playerService.GetScoreboard(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scoreboard})
}
