package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/pkg/helper"
	integration "sugurta/internal/usecase/intgration"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type integrationRoutes struct {
	handlers.BaseHandler
	log                *zap.Logger
	cfg                *config.Config
	enforcer           *casbin.CachedEnforcer
	integrationUsecase integration.Integration
}

func NewIntegrationRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &integrationRoutes{
		log:                option.Logger,
		cfg:                option.Config,
		enforcer:           option.Enforcer,
		integrationUsecase: option.Integration,
	}

	integration := apiV1Group.Group("/integration")
	{
		integration.POST("/create", r.CreateIntegration)
		integration.PUT("/update/:id", r.UpdateIntegration)
		integration.DELETE("/delete/:id", r.DeleteIntegration)
		integration.GET("/owner/:business_id", r.GetIntegrationByBusinessId)
		integration.PUT("/status", r.UpdateStatus)
		integration.GET("/usage/:business_id", r.GetTokenUsageList)
		integration.GET("/existence", r.CheckIntegrationExistence)

	}
}

// CreateIntegration godoc
// @Summary Create a new integration
// @Description Create a new integration with given token and type
// @Tags INTEGRATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param integration body entity.IntegrationCreate true "Integration data"
// @Success 201 {object} status_http.Response{data=string} "Created"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /integration/create [post]
func (i *integrationRoutes) CreateIntegration(c *gin.Context) {
	var req entity.IntegrationCreate
	// bussnesId,code:=helper.GetBusnessIdFromToken(c, i.Config)
	// if code != 0 {
	// 	i.handleResponse(c, status_http.Unauthorized, "Unauthorized")
	// }

	if err := c.ShouldBindJSON(&req); err != nil {
		i.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}
	//req.BusinessId=bussnesId
	if err := i.integrationUsecase.Create(c, &req); err != nil {
		i.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	i.handleResponse(c, status_http.Created, "Integration created successfully")
}

// UpdateIntegration godoc
// @Summary Update integration token
// @Description Update integration token using ID
// @Tags INTEGRATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Integration ID"
// @Param integration body entity.IntegrationUpdateForSwagger true "Update integration"
// @Success 200 {object} status_http.Response{data=string} "Updated"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /integration/update/{id} [put]
func (i *integrationRoutes) UpdateIntegration(c *gin.Context) {
	var req entity.IntegrationUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		i.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}
	id := c.Param("id")
	req.ID = id
	res, err := i.integrationUsecase.Update(c, &req)
	if err != nil {
		i.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	if req.Token != "" {
		BotStart := entity.BotIntegration{
			Token: req.Token,
			Guid:  res.GUID,
		}
		body, err := json.Marshal(BotStart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
			return
		}
		if res.Itype == "bot" {
			resp, err := http.Post("http://ai-seller-bot:8081/start", "application/json", bytes.NewBuffer(body))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
				return
			}
			var BotResp entity.BotIntegrationResponse
			if err := json.NewDecoder(resp.Body).Decode(&BotResp); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
				return
			}
			if BotResp.Code != 0 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bot start" + BotResp.Message})
				return
			}
			i.handleResponse(c, status_http.OK, "Integration updated successfully")
			return
		}
	}
	i.handleResponse(c, status_http.OK, "Integration updated successfully")
}

// UpdateIntegrationStatus godoc
// @Summary Update integration status
// @Description Update integration status using ID
// @Tags INTEGRATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param integration body entity.IntegrationUpdateStatus true "Update integration status"
// @Success 200 {object} status_http.Response{data=string} "Updated"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /integration/status [put]
func (i *integrationRoutes) UpdateStatus(c *gin.Context) {
	var req entity.IntegrationUpdateStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		i.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	res, err := i.integrationUsecase.UpdateStatus(c, &req)
	if err != nil {
		i.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	if res.IntegrationType == "bot" {

		botURL := "http://ai-seller-bot:8081/"
		if req.Status == "active" {
			botURL += "start"
		} else if req.Status == "stop" {
			botURL += "stop"
		} else {
			i.handleResponse(c, status_http.BadRequest, "Invalid status value")
			return
		}

		botReq := entity.BotIntegration{
			Guid:  res.BusinessId,
			Token: res.IntegrationToken,
		}
		if req.Status == "stop" {

			botReq.Token = ""
		}

		body, err := json.Marshal(botReq)
		if err != nil {
			i.handleResponse(c, status_http.InternalServerError, "Failed to marshal bot request: "+err.Error())
			return
		}

		resp, err := http.Post(botURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			i.handleResponse(c, status_http.InternalServerError, "Failed to send request to bot: "+err.Error())
			return
		}
		defer resp.Body.Close()

		var botResp entity.BotIntegrationResponse
		if err := json.NewDecoder(resp.Body).Decode(&botResp); err != nil {
			i.handleResponse(c, status_http.InternalServerError, "Failed to decode bot response: "+err.Error())
			return
		}
		if botResp.Code != 0 {
			i.handleResponse(c, status_http.InternalServerError, "Bot error: "+botResp.Message)
			return
		}
	}

	i.handleResponse(c, status_http.OK, "Integration updated successfully")
}

// DeleteIntegration godoc
// @Summary Delete integration
// @Description Delete integration by ID
// @Tags INTEGRATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Integration ID"
// @Success 200 {object} status_http.Response{data=string} "Deleted"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /integration/delete/{id} [delete]
func (i *integrationRoutes) DeleteIntegration(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		i.handleResponse(c, status_http.BadRequest, "id is required")
		return
	}
	if err := i.integrationUsecase.Delete(c, id); err != nil {
		i.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	i.handleResponse(c, status_http.OK, "Integration deleted successfully")
}

// GetIntegrationByBusinessId godoc
// @Summary Get integration by owner ID
// @Description Retrieve integration data by owner ID
// @Tags INTEGRATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param business_id path string true "Owner ID"
// @Success 200 {object} status_http.Response{data=entity.IntegrationGetResponse} "OK"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /integration/owner/{business_id} [get]
func (i *integrationRoutes) GetIntegrationByBusinessId(c *gin.Context) {
	ownerID := c.Param("business_id")
	if ownerID == "" {
		i.handleResponse(c, status_http.BadRequest, "business_id is required")
		return
	}
	req := &entity.IntegrationRequest{BusinessId: ownerID}
	resp, err := i.integrationUsecase.GetByOwnerID(c, req)
	if err != nil {
		i.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	i.handleResponse(c, status_http.OK, resp)
}

// UpdateIntegrationStatus godoc
// @Summary Update integration status
// @Description Update integration status using ID
// @Tags INTEGRATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param integration body entity.IntegrationUpdateStatus true "Update integration status"
// @Success 200 {object} status_http.Response{data=string} "Updated"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /integration/status [put]
func (i *integrationRoutes) SendTelegramCode(c *gin.Context) {
	var req entity.PhoneNumber
	if err := c.ShouldBindJSON(&req); err != nil {
		i.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	botURL := "http://ai-seller-bot:8081/telegram/send-code"

	body, err := json.Marshal(req)
	if err != nil {
		i.handleResponse(c, status_http.InternalServerError, "Failed to marshal bot request: "+err.Error())
		return
	}

	resp, err := http.Post(botURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		i.handleResponse(c, status_http.InternalServerError, "Failed to send request to bot: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var botResp entity.BotIntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&botResp); err != nil {
		i.handleResponse(c, status_http.InternalServerError, "Failed to decode bot response: "+err.Error())
		return
	}
	if botResp.Code != 0 {
		i.handleResponse(c, status_http.InternalServerError, "Bot error: "+botResp.Message)
		return
	}

	i.handleResponse(c, status_http.OK, "Integration updated successfully")
}

// GetTokenUsageList godoc
// @Summary Get token usage list by business ID
// @Description Retrieves list of token usages for given business with optional filters
// @Tags INTEGRATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param business_id path string true "Business ID"
// @Param source_type query string false "Source type (e.g. bot, telegram)"
// @Param from_date query string false "Start date (RFC3339 format)"
// @Param to_date query string false "End date (RFC3339 format)"
// @Success 200 {object} status_http.Response{data=entity.IntegrationListResponse} "OK"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /integration/usage/{business_id} [get]
func (i *integrationRoutes) GetTokenUsageList(c *gin.Context) {
	businessID := c.Param("business_id")
	if businessID == "" {
		i.handleResponse(c, status_http.BadRequest, "business_id is required")
		return
	}
	userAgent := c.Query("User-Agent")
	fmt.Println(DetectDevice(userAgent))

	var (
		fromDate time.Time
		toDate   time.Time
		err      error
	)

	if from := c.Query("from_date"); from != "" {
		fromDate, err = time.Parse(time.RFC3339, from)
		if err != nil {
			i.handleResponse(c, status_http.BadRequest, "invalid from_date format, use RFC3339")
			return
		}
	}

	if to := c.Query("to_date"); to != "" {
		toDate, err = time.Parse(time.RFC3339, to)
		if err != nil {
			i.handleResponse(c, status_http.BadRequest, "invalid to_date format, use RFC3339")
			return
		}
	}

	req := &entity.IntegrationListRequest{
		BusinessID: businessID,
		SourceType: c.Query("source_type"),
		FromDate:   fromDate,
		ToDate:     toDate,
	}

	resp, err := i.integrationUsecase.GetTokenUsageList(c, req)
	if err != nil {
		i.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	i.handleResponse(c, status_http.OK, resp)
}

// CheckIntegrationExistence godoc
// @Summary Check integration existence by business ID
// @Description Returns which integrations (Telegram account, Telegram bot, Instagram) exist for a specific business
// @Tags INTEGRATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} status_http.Response{data=entity.IntegrationExistenceResponse} "Integrations found"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
// @Router /integration/existence [get]
func (i *integrationRoutes) CheckIntegrationExistence(c *gin.Context) {
	businessID, code := helper.GetBusnessIdFromToken(c, i.cfg)
	if code != 0 {
		i.handleResponse(c, status_http.Unauthorized, "Unauthorized")
		return
	}

	resp, err := i.integrationUsecase.CheckIntegrationExistence(c, businessID)
	if err != nil {
		i.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	i.handleResponse(c, status_http.OK, resp)
}

func DetectDevice(userAgent string) string {
	userAgent = strings.ToLower(userAgent)

	if strings.Contains(userAgent, "mobile") || strings.Contains(userAgent, "android") || strings.Contains(userAgent, "iphone") {
		return "mobile"
	}
	return "desktop"
}

func (h *integrationRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
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
