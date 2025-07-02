package v1

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/pkg/helper"

	"sugurta/internal/usecase/order"
	"sugurta/internal/usecase/settings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// OrderRoutes represents the order controller
type OrderRoutes struct {
	handlers.BaseHandler
	log            *zap.Logger
	cfg            *config.Config
	enforcer       *casbin.CachedEnforcer
	OrderUseCase   order.Order
	SettingsUScase settings.SettingsStorage
}

// NewOrderRoutes creates a new order routes controller
func NewOrderRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &OrderRoutes{
		log:            option.Logger,
		cfg:            option.Config,
		enforcer:       option.Enforcer,
		OrderUseCase:   option.Order,
		SettingsUScase: option.Settings,
	}

	r.Cache = option.Cache
	r.Config = option.Config

	orderGroup := apiV1Group.Group("/orders")
	{
		orderGroup.POST("/create", r.createOrder)
		orderGroup.GET("/get/:id", r.getOrder)
		orderGroup.PUT("/update/:id", r.updateOrder)
		orderGroup.DELETE("/delete/:id", r.deleteOrder)
		orderGroup.GET("/export", r.exportOrders)
		orderGroup.GET("/list", r.ListOrders)
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
	if _, err := uuid.Parse(id); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid UUID format")
		return
	}
	order, err := r.OrderUseCase.Get(c, id)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	r.handleResponse(c, status_http.OK, order)
}

// @Router /orders/list [get]
// @Summary List Orders
// @Description List orders with specific statuses
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param day query int false "Day" default(7)
// @Param status query string true "Status"
// @Param client_id query string false "Client ID"
// @Param search query string false "Search"
// @Param platform query string false "Platform"
// @Param payment_method query string false "Payment Method"
// @Success 200 {object} entity.GetAllOrdersResponse
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) ListOrders(c *gin.Context) {
	allowedStatuses := []string{
		"yangi",
		"yetkazish_kerak",
		"yetkazildi",
		"online_tolov_tasdigi",
		"tolov_qilmoqchi",
		"bekor_qilindi",
		"arxiv",
		"olishga_tayyor",
	}

	// --- Query params
	limit := parseQueryInt(c, "limit", 10)
	offset := parseQueryInt(c, "offset", 0)
	day := parseQueryInt(c, "day", 7)
	status := c.Query("status")

	if !slices.Contains(allowedStatuses, status) {
		r.handleResponse(c, status_http.BadRequest, "Invalid status")
		return
	}

	businessID, code := helper.GetBusnessIdFromToken(c, r.Config)
	if code != 0 {
		r.handleResponse(c, status_http.Unauthorized, "Unauthorized")
		return
	}

	// --- Filter
	filter := &entity.OrderFilter{
		ClientID:      c.Query("client_id"),
		BusinessID:    businessID,
		Status:        status,
		PaymentMethod: c.Query("payment_method"),
		Platform:      c.Query("platform"),
		Search:        c.Query("search"),
		Daye:          day,
	}

	result, err := r.OrderUseCase.List(c, filter, uint64(limit), uint64(offset))
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, result)
}

// @Router /orders/export [get]
// @Summary Export Orders to Excel
// @Description Export orders with applied filters in Excel format
// @Tags Orders
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Param client_id query string false "Client ID"
// @Param business_id query string false "Business ID"
// @Param status query string false "Status"
// @Param search query string false "Search"
// @Param platform query string false "Platform"
// @Param payment_method query string false "Payment Method"
// @Success 200 {file} file "Excel file"
// @Failure 500 {object} status_http.Response{data=string}
func (r *OrderRoutes) exportOrders(c *gin.Context) {
	BusinessID, code := helper.GetBusnessIdFromToken(c, r.Config)
	if code != 0 {
		r.handleResponse(c, status_http.Unauthorized, "Unauthorized")
	}

	filter := &entity.OrderFilter{
		ClientID:      c.Query("client_id"),
		BusinessID:    BusinessID,
		Status:        c.Query("status"),
		PaymentMethod: c.Query("payment_method"),
		Platform:      c.Query("platform"),
		Search:        c.Query("search"),
	}

	fileBytes, err := r.ExportToExcel(c, filter)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename=orders.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
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
	if _, err := uuid.Parse(id); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid UUID format")
		return
	}
	var req entity.OrderUpdate
	BusinessID, code := helper.GetBusnessIdFromToken(c, r.Config)
	if code != 0 {
		r.handleResponse(c, status_http.Unauthorized, "Unauthorized")
	}
	req.BussnesId = BusinessID
	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid request: "+err.Error())
		return
	}

	switch req.Status {
	case "tolov_qilindi":
		req.Status = "yetkazish_kerak"
		req.StatusNumber = 4
	case "yetkazip_berildi":
		req.Status = "yetkazildi"
		req.StatusNumber = 5
	case "bekor_qilindi":
		req.Status = "bekor_qilindi"
		req.StatusNumber = 6
	default:
		r.handleResponse(c, status_http.BadRequest, "Noto‘g‘ri status qiymati")
		return
	}

	statusID, err := r.SettingsUScase.GetStatusByName(c, req.Status, req.BussnesId)
	if err != nil {
		r.handleResponse(c, status_http.BadRequest, "Failed to get status ID: "+err.Error())
		return
	}

	req.StatusID = statusID
	req.ID = id

	// Update order
	if err := r.OrderUseCase.Update(c, &req); err != nil {
		r.handleResponse(c, status_http.InternalServerError, "Failed to update order: "+err.Error())
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
	if _, err := uuid.Parse(id); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid UUID format")
		return
	}
	err := r.OrderUseCase.Delete(c, id)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, "Order deleted successfully")
}

// @Router /orders/products/{id} [get]
// @Summary Get Products By Order ID
// @Description Get list of products related to a specific order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {array} entity.OrderProductBuOrderID
// @Failure 404 {object} status_http.Response{data=string} "Order Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *OrderRoutes) getOrderProductsByOrderID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid UUID format")
		return
	}
	products, err := r.OrderUseCase.GetProductsByOrderID(c, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.handleResponse(c, status_http.NotFound, "Order Not Found")
			return
		}
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, products)
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

func (uc *OrderRoutes) ExportToExcel(ctx context.Context, filter *entity.OrderFilter) ([]byte, error) {
	orders, err := uc.OrderUseCase.List(ctx, filter, 0, 0)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Orders"
	f.SetSheetName("Sheet1", sheet)

	// Header row
	headers := []string{
		"Order ID", "Client ID", "Business ID", "Status", "Payment Method", "Platform", "Total Price", "Created At",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Data rows
	for idx, order := range orders.Items {
		values := []interface{}{
			order.ID,
			order.Client,
			order.BusinessID,
			order.Status,
			order.PaymentMethod,
			order.Platform,
			order.TotalPrice,
			order.CreatedAt,
		}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, idx+2)
			f.SetCellValue(sheet, cell, val)
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parseQueryInt(c *gin.Context, key string, defaultVal int) int {
	valStr := c.DefaultQuery(key, fmt.Sprintf("%d", defaultVal))
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
