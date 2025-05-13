package v1

import (
	"fmt"
	"strconv"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"

	"sugurta/internal/usecase/order"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OrderRoutes represents the order controller
type OrderRoutes struct {
	handlers.BaseHandler
	log          *zap.Logger
	cfg          *config.Config
	enforcer     *casbin.CachedEnforcer
	OrderUseCase order.Order
}

// NewOrderRoutes creates a new order routes controller
func NewOrderRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &OrderRoutes{
		log:          option.Logger,
		cfg:          option.Config,
		enforcer:     option.Enforcer,
		OrderUseCase: option.Order,
	}

	r.Cache = option.Cache
	r.Config = option.Config

	orderGroup := apiV1Group.Group("/orders")
	{
		orderGroup.POST("/create", r.createOrder)
		orderGroup.GET("/get/:id", r.getOrder)
		orderGroup.GET("/list", r.listOrders)
		orderGroup.PUT("/update/:id", r.updateOrder)
		orderGroup.DELETE("/delete/:id", r.deleteOrder)
	}
}

// @Router /orders/create [post]
// @Summary Create Order
// @Description Creates a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body entity.CreateOrderRequest true "CreateOrderRequest"
// @Success 201 {object} status_http.Response
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) createOrder(c *gin.Context) {
	var req entity.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	id, err := r.OrderUseCase.Create(c, &req)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.Created, gin.H{"order_id": id})
}

// @Router /orders/get/{id} [get]
// @Summary Get Order
// @Description Fetch a specific order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} entity.Order
// @Failure 404 {object} status_http.Response{data=string} "Order Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) getOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := r.OrderUseCase.Get(c, id)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	r.handleResponse(c, status_http.OK, order)
}

// @Router /orders/list [get]
// @Summary List Orders
// @Description List all orders with optional filters
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param client_id query string false "Client ID"
// @Param business_id query string false "Business ID"
// @Param status query string false "Status"
// @Param search query string false "Search"
// @Param platform query string false "platform"
// @Param payment_method query string false "Payment Method"
// @Success 200 {object} entity.GetAllOrdersResponse
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) listOrders(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	filter := &entity.OrderFilter{
		ClientID:      c.Query("client_id"),
		BusinessID:    c.Query("business_id"),
		Status:        c.Query("status"),
		PaymentMethod: c.Query("payment_method"),
		Platform:      c.Query("platform"),
		Search:        c.Query("search"),
	}

	result, err := r.OrderUseCase.List(c, filter, uint64(limit), uint64(offset))
	if err != nil {
		fmt.Println(err)
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, result)
}

// @Router /orders/update/{id} [put]
// @Summary Update Order
// @Description Update an existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param order body entity.OrderUpdateForSwagger true "Order Data"
// @Success 200 {object} status_http.Response
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) updateOrder(c *gin.Context) {
	id := c.Param("id")
	var req entity.OrderUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	req.ID = id
	if err := r.OrderUseCase.Update(c, &req); err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, "Order updated successfully")
}

// @Router /orders/delete/{id} [delete]
// @Summary Delete Order
// @Description Delete a specific order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} status_http.Response
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) deleteOrder(c *gin.Context) {
	id := c.Param("id")
	err := r.OrderUseCase.Delete(c, id)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, "Order deleted successfully")
}

func (h *OrderRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
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
