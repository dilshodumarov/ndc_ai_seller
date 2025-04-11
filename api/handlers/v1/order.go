package v1

import (
	"fmt"
	"strconv"
	"strings"
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
		orderGroup.POST("/", r.createOrder)
		orderGroup.GET("/:id", r.getOrder)
		orderGroup.GET("/", r.listOrders)
		orderGroup.PUT("/:id", r.updateOrder)
		orderGroup.DELETE("/:id", r.deleteOrder)
	}
}

// @Router /orders [post]
// @Summary Create Order
// @Description Creates a new order
// @Tags Orders
// @Accept json
// @Produce json
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

// @Router /orders/{id} [get]
// @Summary Get Order
// @Description Fetch a specific order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} entity.Order
// @Failure 404 {object} status_http.Response{data=string} "Order Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) getOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := r.OrderUseCase.Get(c, map[string]string{"id": id})
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	r.handleResponse(c, status_http.OK, order)
}

// @Router /orders [get]
// @Summary List Orders
// @Description List all orders with optional filters
// @Tags Orders
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param filter query string false "Filter by id, client_id, integration_id, status"
// @Success 200 {object} entity.GetAllOrdersResponse
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) listOrders(c *gin.Context) {
	
	// Extract limit and offset from the query parameters
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	filterQuery := c.DefaultQuery("filter", "")

	// Convert limit and offset to integers (or use defaults)
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10 // default value
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0 // default value
	}

	// Convert filter query into a map (assuming filters are passed as comma-separated key-value pairs)
	filter := make(map[string]string)
	if filterQuery != "" {
		for _, pair := range strings.Split(filterQuery, ",") {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				filter[parts[0]] = parts[1]
			}
		}
	}
	
	// Call the List method with the constructed filter map
	fmt.Println("order uscase: ",r.OrderUseCase)
	orders, err := r.OrderUseCase.List(c, uint64(limitInt), uint64(offsetInt), filter)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		
		return
	}
	
	r.handleResponse(c, status_http.OK, orders)
}

// @Router /orders/{id} [put]
// @Summary Update Order
// @Description Update an existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param order body entity.OrderUpdate true "Order Data"
// @Success 200 {object} status_http.Response
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) updateOrder(c *gin.Context) {
	id := c.Param("id")
	var order entity.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		r.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	order.ID = id
	err := r.OrderUseCase.Update(c, &order)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, "Order updated successfully")
}

// @Router /orders/{id} [delete]
// @Summary Delete Order
// @Description Delete a specific order
// @Tags Orders
// @Accept json
// @Produce json
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
