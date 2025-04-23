package v1

import (
	"sugurta/api/handlers"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type telegramRoutes struct {
	handlers.BaseHandler
	log *zap.Logger
	cfg *config.Config
}

// NewTelegramRoutes registers Telegram-related routes
func NewTelegramRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &telegramRoutes{
		log: option.Logger,
		cfg: option.Config,
	}

	telegram := apiV1Group.Group("/telegram")
	{
		telegram.POST("/send-code", r.SendTelegramCode)
		telegram.POST("/verify", r.SendTelegramVerify)
	}
}

// SendTelegramCode godoc
// @Summary Send phone number to Telegram bot
// @Description Sends a phone number to the Python backend to receive a code
// @Tags TELEGRAM
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param phone body entity.PhoneNumber true "Phone number"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /telegram/send-code [post]
func (h *telegramRoutes) SendTelegramCode(c *gin.Context) {
	var phone entity.PhoneNumber
	if err := c.ShouldBindJSON(&phone); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := SendTelegramCode(phone)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

// SendTelegramVerify godoc
// @Summary Verify Telegram code
// @Description Verifies the code (and optional password) received from the user
// @Tags TELEGRAM
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param verify body entity.CodeInput true "Verification input"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /telegram/verify [post]
func (h *telegramRoutes) SendTelegramVerify(c *gin.Context) {
	var input entity.CodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := SendTelegramVerify(input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}
