package v1

import (
	"strconv"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/pkg/helper"
	"sugurta/internal/usecase/settings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type settingsRoutes struct {
	handlers.BaseHandler
	log        *zap.Logger
	cfg        *config.Config
	enforcer   *casbin.CachedEnforcer
	settingsUC settings.SettingsStorage
}

func NewSettingsRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &settingsRoutes{
		log:        option.Logger,
		cfg:        option.Config,
		enforcer:   option.Enforcer,
		settingsUC: option.Settings,
	}

	settings := apiV1Group.Group("/settings/order-status")
	{
		//settings.POST("/create", r.CreateOrderStatus)
		settings.GET("/get/:guid", r.GetOrderStatus)
		settings.PUT("/update/:guid", r.UpdateOrderStatus)
		settings.DELETE("/delete/:guid", r.DeleteOrderStatus)
		settings.GET("/list", r.ListOrderStatus)
	}
	settingsai := apiV1Group.Group("settings/ai")
	{
		//settingsai.POST("/create", r.CreateSettings)
		settingsai.GET("/get/:guid", r.GetSettings)
		settingsai.PUT("/update", r.UpdateSettings)
		settingsai.DELETE("/delete/:guid", r.DeleteSettings)
		settingsai.GET("/list", r.ListSettingsByBusinessID)
		settingsai.GET("/by-businessid", r.GetSettingsByName)
		settingsai.PUT("/prompt-orders/:guid", r.UpdatePromptOrders)
		settingsai.GET("/prompt-orders", r.GetPromptOrders)
	}

}

// // CreateOrderStatus godoc
// // @Summary Create a new order status
// // @Description Creates a new order status for the given business
// // @Tags SETTINGS
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param order_status body entity.CreateOrderStatusRequestForswagger true "Order Status"
// // @Success 201 {object} status_http.Response{data=string} "Created"
// // @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// // @Failure 500 {object} status_http.Response{data=string} "Server Error"
// // @Router /settings/order-status/create [post]
// func (r *settingsRoutes) CreateOrderStatus(c *gin.Context) {
// 	var req entity.CreateOrderStatusRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		r.handleResponse(c, status_http.BadRequest, err.Error())
// 		return
// 	}
// 	BusinessID, code := helper.GetBusnessIdFromToken(c, r.cfg)
// 	if code != 0 {
// 		r.handleResponse(c, status_http.Unauthorized, "Unauthorized")
// 	}
// 	req.BusinessID = BusinessID
// 	if err := r.settingsUC.Create(c, &req); err != nil {
// 		r.handleResponse(c, status_http.InternalServerError, err.Error())
// 		return
// 	}
// 	r.handleResponse(c, status_http.Created, "Order status created successfully")
// }

// GetOrderStatus godoc
// @Summary Get a specific order status
// @Description Retrieves order status details by GUID
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param guid path string true "Order Status GUID"
// @Success 200 {object} status_http.Response{data=entity.OrderStatus} "OK"
// @Failure 404 {object} status_http.Response{data=string} "Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/order-status/get/{guid} [get]
func (r *settingsRoutes) GetOrderStatus(c *gin.Context) {
	guid := c.Param("guid")

	res, err := r.settingsUC.Get(c, guid)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	r.handleResponse(c, status_http.OK, res)
}

// UpdateOrderStatus godoc
// @Summary Update an order status
// @Description Updates an existing order status by GUID
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param guid path string true "Order Status GUID"
// @Param order_status body entity.UpdateOrderStatusRequest true "Order Status Update"
// @Success 200 {object} status_http.Response{data=string} "Updated"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/order-status/update/{guid} [put]
func (r *settingsRoutes) UpdateOrderStatus(c *gin.Context) {
	guid := c.Param("guid")

	var req entity.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}
	req.GUID = guid
	if err := r.settingsUC.Update(c, &req); err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	r.handleResponse(c, status_http.OK, "Order status updated successfully")
}

// DeleteOrderStatus godoc
// @Summary Delete an order status
// @Description Deletes an order status by GUID
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param guid path string true "Order Status GUID"
// @Success 200 {object} status_http.Response{data=string} "Deleted"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/order-status/delete/{guid} [delete]
func (r *settingsRoutes) DeleteOrderStatus(c *gin.Context) {
	guid := c.Param("guid")
	if err := r.settingsUC.Delete(c, guid); err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	r.handleResponse(c, status_http.OK, "Order status deleted successfully")
}

// ListOrderStatus godoc
// @Summary List all order statuses for a business
// @Description Returns a list of order statuses for a specific business
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "status name"
// @Success 200 {object} status_http.Response{data=[]entity.OrderStatus} "OK"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/order-status/list [get]
func (r *settingsRoutes) ListOrderStatus(c *gin.Context) {
	Status:=c.Query("status")
	BusinessID, code := helper.GetBusnessIdFromToken(c, r.cfg)
	if code != 0 {
		r.handleResponse(c, status_http.Unauthorized, "Unauthorized")
	}
	res, err := r.settingsUC.List(c, entity.OrderStatusFilter{
		BusinessID: BusinessID,
		Status: Status,
	})
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	r.handleResponse(c, status_http.OK, res)
}

// // CreateSettings godoc
// // @Summary Create new settings
// // @Description Creates new settings entry in the database
// // @Tags SETTINGS
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param settings body entity.CreateSettingsRequest true "Settings create request"
// // @Success 201 {string} string "Created successfully"
// // @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// // @Failure 500 {object} status_http.Response{data=string} "Server Error"
// // @Router /settings/ai/create [post]
// func (r *settingsRoutes) CreateSettings(c *gin.Context) {
// 	var req entity.CreateSettingsRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		r.handleResponse(c, status_http.BadRequest, err.Error())
// 		return
// 	}
// 	BusinessID, code := helper.GetBusnessIdFromToken(c, r.cfg)
// 	if code != 0 {
// 		r.handleResponse(c, status_http.Unauthorized, "Unauthorized")
// 	}
// 	req.BusinessID = BusinessID

// 	if err := r.settingsUC.CreateSettings(c, &req); err != nil {
// 		r.handleResponse(c, status_http.InternalServerError, err.Error())
// 		return
// 	}

// 	r.handleResponse(c, status_http.Created, "Created successfully")
// }

// GetSettings godoc
// @Summary Get settings by GUID
// @Description Retrieves settings by GUID
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param guid path string true "Settings GUID"
// @Success 200 {object} status_http.Response{data=entity.Settings} "OK"
// @Failure 404 {object} status_http.Response{data=string} "Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/ai/{guid} [get]
func (r *settingsRoutes) GetSettings(c *gin.Context) {
	guid := c.Param("guid")
	res, err := r.settingsUC.GetSettings(c, guid)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	if res == nil {
		r.handleResponse(c, status_http.NotFound, "Settings not found")
		return
	}
	r.handleResponse(c, status_http.OK, res)
}

// UpdateSettings godoc
// @Summary Update existing settings
// @Description Updates settings entry by GUID
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param settings body entity.UpdateSettingsRequest true "Settings update request"
// @Success 200 {string} string "Updated successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 404 {object} status_http.Response{data=string} "Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/ai/update [put]
func (r *settingsRoutes) UpdateSettings(c *gin.Context) {
	var req entity.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	err := r.settingsUC.UpdateSettings(c, &req)
	if err != nil {
		if err.Error() == "no rows affected" {
			r.handleResponse(c, status_http.NotFound, "Settings not found")
			return
		}
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, "Updated successfully")
}

// DeleteSettings godoc
// @Summary Delete settings by GUID
// @Description Deletes settings entry by GUID
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param guid path string true "Settings GUID"
// @Success 200 {string} string "Deleted successfully"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/ai/delete/{guid} [delete]
func (r *settingsRoutes) DeleteSettings(c *gin.Context) {
	guid := c.Param("guid")
	err := r.settingsUC.DeleteSettings(c, guid)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, "Deleted successfully")
}

// ListSettingsByBusinessID godoc
// @Summary List settings by business ID
// @Description Retrieves a list of settings filtered by business ID
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Settings "OK"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/ai/list [get]
func (r *settingsRoutes) ListSettingsByBusinessID(c *gin.Context) {
	BusinessID, code := helper.GetBusnessIdFromToken(c, r.cfg)
	if code != 0 {
		r.handleResponse(c, status_http.Unauthorized, "Unauthorized")
	}
	settings, err := r.settingsUC.ListSettingsByBusinessID(c, BusinessID)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	r.handleResponse(c, status_http.OK, settings)
}

// GetSettingsByBussnesID godoc
// @Summary Get settings by  business ID
// @Description Retrieves settings by  business ID
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.Settings "OK"
// @Failure 404 {object} status_http.Response{data=string} "Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/ai/by-businessid [get]
func (r *settingsRoutes) GetSettingsByName(c *gin.Context) {
	BusinessID, code := helper.GetBusnessIdFromToken(c, r.cfg)
	if code != 0 {
		r.handleResponse(c, status_http.Unauthorized, "Unauthorized")
	}
	res, err := r.settingsUC.GetSettingsBussnesId(c, BusinessID)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	if res == nil {
		r.handleResponse(c, status_http.NotFound, "Settings not found")
		return
	}

	r.handleResponse(c, status_http.OK, res)
}

// UpdatePromptOrders godoc
// @Summary Update only the prompt_orders field
// @Description Updates the prompt_orders JSONB field of a specific settings entry by GUID
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param guid path string true "Settings GUID"
// @Param prompt_orders body entity.UpdatePromptOrdersRequest true "Prompt orders map"
// @Success 200 {string} string "Prompt orders updated successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 404 {object} status_http.Response{data=string} "Settings not found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/ai/prompt-orders/{guid} [put]
func (r *settingsRoutes) UpdatePromptOrders(c *gin.Context) {
	guid := c.Param("guid")
	if guid == "" {
		r.handleResponse(c, status_http.BadRequest, "GUID is required")
		return
	}

	var req entity.UpdatePromptOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	promptOrderMap := make(map[string]string)

	for _, status := range req.OrderStatus {
		if status.GUID != "" {
			err := r.settingsUC.Update(c, &status)
			if err != nil {
				if err.Error() == "no rows affected" {
					r.handleResponse(c, status_http.NotFound, "Settings not found")
					return
				}
				r.handleResponse(c, status_http.InternalServerError, err.Error())
				return
			}
		}

		// Agar promt_number ixtiyoriy >0 bo‘lsa va promt bo‘sh bo‘lmasa — mapga qo‘shamiz
		if status.PromtNumber > 0 && status.Promt != "" {
			key := strconv.Itoa(status.PromtNumber)
			promptOrderMap[key] = status.Promt
		}
	}

	err := r.settingsUC.UpdatePromptOrders(c, guid, promptOrderMap)
	if err != nil {
		if err.Error() == "no rows affected" {
			r.handleResponse(c, status_http.NotFound, "Settings not found")
			return
		}
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, "Prompt orders updated successfully")
}

// GetPromptOrders godoc
// @Summary Get prompt orders
// @Description Returns prompt orders (2, 3, 4, 6) from JSONB field
// @Tags SETTINGS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.GetPromptOrdersResponse
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 404 {object} status_http.Response{data=string} "Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /settings/ai/prompt-orders [get]
func (r *settingsRoutes) GetPromptOrders(c *gin.Context) {
	BusinessID, code := helper.GetBusnessIdFromToken(c, r.cfg)
	if code != 0 {
		r.handleResponse(c, status_http.Unauthorized, "Unauthorized")
		return
	}

	data, err := r.settingsUC.GetPromptOrders(c, BusinessID)
	if err != nil {
		if err.Error() == "not found" {
			r.handleResponse(c, status_http.NotFound, "Settings not found")
			return
		}
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	OrderStatus, err := r.settingsUC.List(c, entity.OrderStatusFilter{BusinessID: BusinessID})
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	// Promptlarni StatusNumberga bog‘lash
	for _, status := range OrderStatus {
		statusNumberStr := strconv.Itoa(status.StatusNumber)
		for _, prompt := range data {
			if prompt.Number == statusNumberStr {
				status.Prompts = prompt
				break
			}
		}
	}

	r.handleResponse(c, status_http.OK, entity.GetPromptOrdersResponse{
		OrderStatus: OrderStatus,
		Id:          data[0].Guid,
	})
}

func (h *settingsRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
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
