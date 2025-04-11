package v1

import (
	"context"
	"time"

	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/pkg/helper"
	clienttype "sugurta/internal/usecase/client-type"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type clientTypeRoutes struct {
	handlers.BaseHandler
	log          *zap.Logger
	cfg          *config.Config
	enforcer     *casbin.CachedEnforcer
	clientUscase clienttype.ClientType
}

// NewAuthRoutes creates a new auth routes controller
func NewClientTypeRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &clientTypeRoutes{
		log:          option.Logger,
		cfg:          option.Config,
		enforcer:     option.Enforcer,
		clientUscase: option.ClientType,
	}

	clientTypeGroup := apiV1Group.Group("/client-type")
	{
		clientTypeGroup.POST("/", r.createClientType)
		clientTypeGroup.GET("/:id", r.getClientTypeByID)
		clientTypeGroup.GET("/", r.getClientTypeList)
		clientTypeGroup.PUT("/", r.updateClientType)
		clientTypeGroup.DELETE("/:id", r.deleteClientType)
	}
}

// @Router /client-type [post]
// @Summary Create ClientType
// @Description Create a new ClientType
// @Tags CLIENT TYPE
// @Accept json
// @Produce json
// @Param data body entity.CreateClientTypeRequest true "Data"
// @Success 201 {object} status_http.Response{data=string} "Success"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (ct *clientTypeRoutes) createClientType(c *gin.Context) {
	var req entity.CreateClientTypeRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		ct.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = ct.clientUscase.Create(ctx, &req)
	if err != nil {
		ct.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	ct.handleResponse(c, status_http.Created, "created successfully")
}

// @Router /client-type/{id} [get]
// @Summary Get ClientType By ID
// @Description Get ClientType by ID
// @Tags CLIENT TYPE
// @Accept json
// @Produce json
// @Param id path string true "ClientType ID"
// @Success 200 {object} status_http.Response{data=entity.ClientType}
// @Failure 400 {object} status_http.Response{data=string} "Bad request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (ct *clientTypeRoutes) getClientTypeByID(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if id == "" {
		ct.handleResponse(c, status_http.BadRequest, "ClientType ID is required")
		return
	}

	resp, err := ct.clientUscase.Get(ctx, id)
	if err != nil {
		ct.handleResponse(c, status_http.NotFound, err.Error())
		return
	}

	ct.handleResponse(c, status_http.OK, resp)
}

// @Router /client-type [get]
// @Summary Get All ClientTypes
// @Description Get List of ClientTypes
// @Tags CLIENT TYPE
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} status_http.Response{data=entity.ClientTypeListResponse} "Data"
// @Failure 400 {object} status_http.Response{data=string} "Bad request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (ct *clientTypeRoutes) getClientTypeList(c *gin.Context) {
	page, limit := helper.GetPaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := ct.clientUscase.List(ctx, page, limit)
	if err != nil {
		ct.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	ct.handleResponse(c, status_http.OK, resp)
}

// @Router /client-type [put]
// @Summary Update ClientType
// @Description Update ClientType by ID
// @Tags CLIENT TYPE
// @Accept json
// @Produce json
// @Param data body entity.UpdateClientType true "Data"
// @Success 200 {object} status_http.Response{data=string} "Success"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (ct *clientTypeRoutes) updateClientType(c *gin.Context) {

	var req entity.UpdateClientType

	err := c.ShouldBindJSON(&req)
	if err != nil {
		ct.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = ct.clientUscase.Update(ctx, &req)
	if err != nil {
		ct.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	ct.handleResponse(c, status_http.OK, "Successfully updated")
}

// @Router /client-type/{id} [delete]
// @Summary Delete ClientType
// @Description Delete ClientType by ID
// @Tags CLIENT TYPE
// @Accept json
// @Produce json
// @Param id path string true "ClientType ID"
// @Success 200 {object} status_http.Response{data=string} "Successfully deleted"
// @Success 400 {object} status_http.Response{data=string} "Bad request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (ct *clientTypeRoutes) deleteClientType(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		ct.handleResponse(c, status_http.BadRequest, "ClientType ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := ct.clientUscase.Delete(ctx, id)
	if err != nil {
		ct.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	ct.handleResponse(c, status_http.OK, "client type deleted successfully")
}

// handleResponse handles the HTTP response
func (h *clientTypeRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
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
