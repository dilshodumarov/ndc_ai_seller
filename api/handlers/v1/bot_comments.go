package v1

import (
	"fmt"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/usecase/bot_comments"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type botCommentsRoutes struct {
	handlers.BaseHandler
	log                *zap.Logger
	cfg                *config.Config
	enforcer           *casbin.CachedEnforcer
	botCommentsUsecase botcomments.BotCommandStorage
}

func NewBotCommentsRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &botCommentsRoutes{
		log:                option.Logger,
		cfg:                option.Config,
		enforcer:           option.Enforcer,
		botCommentsUsecase: option.BotComments,
	}

	botGroup := apiV1Group.Group("/bot-comments")
	{
		botGroup.POST("/create", r.createBotCommand)
		botGroup.GET("/get/:id", r.getBotCommand)
		botGroup.PUT("/update", r.updateBotCommand)
		botGroup.DELETE("/delete/:id", r.deleteBotCommand)
		botGroup.GET("/list/:integration_id", r.listBotCommands)
	}
}

// @Summary Create bot command
// @Description Create a new bot command for integration
// @Tags BOT_COMMENTS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body entity.BotCommandRequest true "Bot Command Data"
// @Success 201 {object} status_http.Response{data=map[string]string} "GUID"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /bot-comments/create [post]
func (h *botCommentsRoutes) createBotCommand(c *gin.Context) {
	var req entity.BotCommandRequest
	fmt.Println(1111111)
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, status_http.BadRequest, "invalid request data")
		return
	}
	
	guid, err := h.botCommentsUsecase.CreateBotCommand(c, req)
	if err != nil {
		h.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, status_http.Created, map[string]string{"guid": guid})
}

// @Summary Get bot command
// @Description Get a bot command by ID (guid)
// @Tags BOT_COMMENTS
// @Security BearerAuth
// @Produce json
// @Param id path string true "Command GUID"
// @Success 200 {object} status_http.Response{data=entity.BotCommandResponse}
// @Failure 404 {object} status_http.Response{data=string} "Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /bot-comments/get/{id} [get]
func (h *botCommentsRoutes) getBotCommand(c *gin.Context) {
	guid := c.Param("id")

	result, err := h.botCommentsUsecase.GetBotCommand(c, guid)
	if err != nil {
		h.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	if result == nil {
		h.handleResponse(c, status_http.NotFound, "bot command not found")
		return
	}

	h.handleResponse(c, status_http.OK, result)
}

// @Summary Update bot command
// @Description Update bot command fields (command, response)
// @Tags BOT_COMMENTS
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body entity.BotCommandUpdateRequest true "Update Data"
// @Success 200 {object} status_http.Response{data=string} "Updated successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /bot-comments/update [put]
func (h *botCommentsRoutes) updateBotCommand(c *gin.Context) {
	var req entity.BotCommandUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, status_http.BadRequest, "invalid request data")
		return
	}

	if err := h.botCommentsUsecase.UpdateBotCommand(c, req); err != nil {
		h.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, status_http.OK, "bot command updated successfully")
}

// @Summary Delete bot command
// @Description Delete a bot command by GUID
// @Tags BOT_COMMENTS
// @Security BearerAuth
// @Produce json
// @Param id path string true "Command GUID"
// @Success 200 {object} status_http.Response{data=string} "Deleted successfully"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /bot-comments/delete/{id} [delete]
func (h *botCommentsRoutes) deleteBotCommand(c *gin.Context) {
	guid := c.Param("id")

	err := h.botCommentsUsecase.DeleteBotCommand(c, guid)
	if err != nil {
		h.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, status_http.OK, "bot command deleted successfully")
}

// @Summary List bot commands
// @Description List all bot commands for an integration
// @Tags BOT_COMMENTS
// @Security BearerAuth
// @Produce json
// @Param integration_id path string true "Integration ID"
// @Success 200 {object} status_http.Response{data=[]entity.BotCommandResponse}
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /bot-comments/list/{integration_id} [get]
func (h *botCommentsRoutes) listBotCommands(c *gin.Context) {
	integrationID := c.Param("integration_id")

	list, err := h.botCommentsUsecase.ListBotCommands(c, integrationID)
	if err != nil {
		h.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, status_http.OK, list)
}

func (h *botCommentsRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
	switch code := status.Code; {
	case code < 400:
	default:
		h.log.Error(
			"response",
			zap.Int("code", status.Code),
			zap.String("status", status.Status),
			zap.Any("description", status.Description),
			zap.Any("data", data),
			zap.Any("custom_message", status.CustomMessage),
		)
	}

	c.JSON(status.Code, status_http.Response{
		Status:        status.Status,
		Description:   status.Description,
		Data:          data,
		CustomMessage: status.CustomMessage,
	})
}
