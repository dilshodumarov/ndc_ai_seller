package v1

import (
	"fmt"
	"sugurta/api/handlers"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/usecase/telegram"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type telegramRoutes struct {
	handlers.BaseHandler
	log      *zap.Logger
	cfg      *config.Config
	telegram telegram.TelegramAccount
}

// NewTelegramRoutes registers Telegram-related routes
func NewTelegramRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &telegramRoutes{
		log:      option.Logger,
		cfg:      option.Config,
		telegram: option.Telegram,
	}

	telegram := apiV1Group.Group("/telegram")
	{
		telegram.POST("/send-code", r.SendTelegramCode)
		telegram.POST("/verify", r.SendTelegramVerify)
		telegram.POST("/start-session", r.StartTelegramSession)
		telegram.POST("/stop-session", r.StopTelegramSession)

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
// @Success 200 {object} entity.BotIntegrationResponse "Success"
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
	if resp.Code == 0 {
		_, err := h.telegram.Create(c, entity.CreateTelegramAccountRequest{
			Number:     input.Phone,
			BusinessID: input.BussnesId,
			From: "telegram",
		})
		if err != nil {
			fmt.Println(err)
			c.JSON(500, "server error")
			return
		}
		fmt.Println("telegram accaunt yaratildi")
	}

	c.JSON(200, resp)
}

// StartTelegramSession godoc
// @Summary Start Telegram session
// @Description Starts a new Telegram session for the given business ID
// @Tags TELEGRAM
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.PhoneNumber true "Start session request"
// @Success 200 {object} map[string]interface{} "Session started"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /telegram/start-session [post]
func (h *telegramRoutes) StartTelegramSession(c *gin.Context) {
	var input entity.PhoneNumber
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	resp, err := SendTelegramStartSession(input)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, "server error")
		return
	}

	err = h.telegram.Update(c, entity.UpdateTelegramAccountRequest{
		Phone:  input.Phone,
		Status: "start",
	})

	if err != nil {
		c.JSON(500, "server error")
		return
	}
	c.JSON(200, resp)
}

// StopTelegramSession godoc
// @Summary Stop Telegram session
// @Description Stops the Telegram session for the given business ID
// @Tags TELEGRAM
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.PhoneNumber true "Stop session request"
// @Success 200 {object} map[string]interface{} "Session stopped"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /telegram/stop-session [post]
func (h *telegramRoutes) StopTelegramSession(c *gin.Context) {
	var input entity.PhoneNumber
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	resp, err := SendTelegramStopSession(input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = h.telegram.Update(c, entity.UpdateTelegramAccountRequest{
		Phone:  input.Phone,
		Status: "stop",
	})

	if err != nil {
		c.JSON(500, "server error")
		return
	}
	c.JSON(200, resp)
}

// // ListTelegramSessions godoc
// // @Summary List all active Telegram sessions
// // @Description Returns a list of all active Telegram sessions
// // @Tags TELEGRAM
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {array} entity.TelegramSessionResponse "List of sessions"
// // @Failure 500 {object} map[string]interface{} "Internal Server Error"
// // @Router /telegram/sessions [get]
// func (h *telegramRoutes) ListTelegramSessions(c *gin.Context) {
// 	resp, err := ListTelegramSessions()
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(200, resp)
// }
