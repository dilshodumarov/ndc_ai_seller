package v1

import (
	"fmt"
	"strconv"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/usecase/notification"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// notificationRoutes defines the routes for notification management
type notificationRoutes struct {
	handlers.BaseHandler
	notificationUsecase notification.Notification
	log                 *zap.Logger
}

// NewNotificationRoutes creates new notification routes
func NewNotificationRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &notificationRoutes{
		notificationUsecase: option.Notification,
		log:                 option.Logger,
	}

	notification := apiV1Group.Group("/notification")
	{
		notification.POST("/create", r.CreateNotification)
		notification.GET("/get/:id", r.GetNotification)
		notification.PUT("/update/:id", r.UpdateNotification)
		notification.DELETE("/delete/:id", r.DeleteNotification)
		notification.GET("/list", r.ListNotifications)
	}
}

// CreateNotification godoc
// @Summary Create a new notification
// @Description Create a new notification with the provided details
// @Tags NOTIFICATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body entity.CreateNotificationRequest true "Notification details"
// @Success 201 {object} status_http.Response{data=string} "Success"
// @Response 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /notification/create [post]
func (n *notificationRoutes) CreateNotification(c *gin.Context) {
	var notification entity.CreateNotificationRequest

	if err := c.ShouldBindJSON(&notification); err != nil {
		n.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	id, err := n.notificationUsecase.Create(c, &notification)
	if err != nil {
		n.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	n.handleResponse(c, status_http.Created, "Notification created successfully id:  "+id)
}

// GetNotification godoc
// @Summary Get a notification by ID
// @Description Get a notification by its unique identifier
// @Tags NOTIFICATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} entity.GetNotification "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /notification/get/{id} [get]
func (n *notificationRoutes) GetNotification(c *gin.Context) {
	id := c.Param("id")

	notification, err := n.notificationUsecase.Get(c, id)
	if err != nil {
		n.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	if !notification.IsRead{
		err:=n.notificationUsecase.MarkAsRead(c, id)
		if err!=nil{
			fmt.Println(err)
		}
	}
	n.handleResponse(c, status_http.OK, notification)
}

// UpdateNotification godoc
// @Summary Update a notification
// @Description Update an existing notification with new details
// @Tags NOTIFICATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Param notification body entity.UpdateNotificationRequest true "Notification details"
// @Success 200 {object} status_http.Response{data=string} "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /notification/update/{id} [put]
func (n *notificationRoutes) UpdateNotification(c *gin.Context) {
	var notification entity.UpdateNotificationRequest
	id := c.Param("id")

	if err := c.ShouldBindJSON(&notification); err != nil {
		n.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	notification.GUID = id
	if err := n.notificationUsecase.Update(c, &notification); err != nil {
		n.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	n.handleResponse(c, status_http.OK, "Notification updated successfully")
}

// DeleteNotification godoc
// @Summary Delete a notification
// @Description Delete a notification by its unique identifier
// @Tags NOTIFICATION
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} status_http.Response{data=string} "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /notification/delete/{id} [delete]
func (n *notificationRoutes) DeleteNotification(c *gin.Context) {
	id := c.Param("id")

	if err := n.notificationUsecase.Delete(c, id); err != nil {
		n.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	n.handleResponse(c, status_http.OK, "Notification deleted successfully")
}

// ListNotifications godoc
// @Summary List all notifications
// @Description List notifications by user ID with pagination
// @Tags NOTIFICATION
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} entity.ListNotification "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /notification/list [get]
func (n *notificationRoutes) ListNotifications(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := entity.ListNotificationRequest{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	}

	notifications, err := n.notificationUsecase.List(c, req)
	if err != nil {
		fmt.Println(err)
		n.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	n.handleResponse(c, status_http.OK, notifications)
}

func (h *notificationRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
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
