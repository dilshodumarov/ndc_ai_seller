package v1

import (
	"database/sql"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/usecase/category"

	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type categoryRoutes struct {
	handlers.BaseHandler
	log            *zap.Logger
	cfg            *config.Config
	enforcer       *casbin.CachedEnforcer
	categoryUscase category.Category
}

func NewcategoryRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &categoryRoutes{
		log:            option.Logger,
		cfg:            option.Config,
		enforcer:       option.Enforcer,
		categoryUscase: option.Category,
	}

	categoriesGroup := apiV1Group.Group("/category")
	{
		categoriesGroup.POST("", r.createCategory)
		categoriesGroup.GET("/:id", r.GetCategoryByID)
		categoriesGroup.PUT("", r.updateCategory)
		categoriesGroup.DELETE("/:id", r.deleteCategory)
	}
}

// @Router /category [post]
// @Summary Create a new category
// @Description Create a new category in the database
// @Tags CATEGORY
// @Accept json
// @Produce json
// @Param category body entity.CreateCategoryRequest true "Category Details"
// @Success 201 {object} status_http.Response{data=string} "Category created successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (h *categoryRoutes) createCategory(c *gin.Context) {
	var category entity.CreateCategoryRequest
	if err := c.ShouldBindJSON(&c); err != nil {
		h.handleResponse(c, status_http.BadRequest, "invalid request data")
		return
	}

	err := h.categoryUscase.Create(c, &category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})
}

// @Router /category/{id} [get]
// @Summary Get category by ID
// @Description Get category details by ID
// @Tags CATEGORY
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} status_http.Response{data=entity.Category} "Category"
// @Failure 400 {object} status_http.Response{data=string} "Bad request data"
// @Failure 404 {object} status_http.Response{data=string} "Category Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (h *categoryRoutes) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.handleResponse(c, status_http.BadRequest, "id is required")
		return
	}

	cat, err := h.categoryUscase.Get(c, id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.handleResponse(c, status_http.NotFound, "category not found")
			return
		}
		h.handleResponse(c, status_http.InternalServerError, "error getting category")
		return
	}

	c.JSON(http.StatusOK, cat)
}

// @Router /category [put]
// @Summary Update an existing category
// @Description Update category details by ID
// @Tags CATEGORY
// @Accept json
// @Produce json
// @Param category body entity.Category true "Category Details"
// @Success 200 {object} status_http.Response{data=string} "Category updated successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad request data"
// @Failure 404 {object} status_http.Response{data=string} "Category Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (h *categoryRoutes) updateCategory(c *gin.Context) {

	var category entity.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&c); err != nil {
		h.handleResponse(c, status_http.BadRequest, "invalid request data")
		return
	}

	err := h.categoryUscase.Update(c, &category)
	if err != nil {
		if err == sql.ErrNoRows {
			h.handleResponse(c, status_http.NotFound, "category not found")
			return
		}
		h.handleResponse(c, status_http.InternalServerError, "error updating category")
		return
	}

	h.handleResponse(c, status_http.OK, "category updated successfully")
}

// @Router /category/{id} [delete]
// @Summary Delete a category by ID
// @Description Delete a category from the database by ID
// @Tags CATEGORY
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} status_http.Response{data=string} "Category deleted successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request Data"
// @Failure 404 {object} status_http.Response{data=string} "Category Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (h *categoryRoutes) deleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.handleResponse(c, status_http.BadRequest, "id is required")
		return
	}

	err := h.categoryUscase.Delete(c, id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.handleResponse(c, status_http.NotFound, "category not found")
			return
		}
		h.handleResponse(c, status_http.InternalServerError, "error deleting category")
		return
	}

	h.handleResponse(c, status_http.OK, "category deleted successfully")
}

func (h *categoryRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
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
