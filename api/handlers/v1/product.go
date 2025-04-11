package v1

import (
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/usecase/product"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type productRoutes struct {
	handlers.BaseHandler
	log           *zap.Logger
	cfg           *config.Config
	enforcer      *casbin.CachedEnforcer
	productUscase product.Product
}

func NewProductRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &productRoutes{
		log:           option.Logger,
		cfg:           option.Config,
		enforcer:      option.Enforcer,
		productUscase: option.Product,
	}

	productGroup := apiV1Group.Group("/product")
	{
		productGroup.POST("", r.createProduct)
		productGroup.GET("/:id", r.getProductByID)
		productGroup.PUT("", r.updateProduct)
		productGroup.DELETE("/:id", r.deleteProduct)
	}

}

// @Router /product [post]
// @Summary Create a new product
// @Description Create a new product in the database
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param product body entity.CreateProductRequest true "Product Details"
// @Success 201 {object} status_http.Response{data=string} "Product created successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (p *productRoutes) createProduct(c *gin.Context) {
	var req entity.CreateProductRequest
	if err := c.ShouldBindJSON(&p); err != nil {
		p.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	err := p.productUscase.Create(c, &req)
	if err != nil {
		p.handleResponse(c, status_http.InternalServerError, "error while creating product")
		return
	}

	p.handleResponse(c, status_http.Created, "Product created successfully")
}

// @Router /product/{id} [get]
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} status_http.Response{data=string} "Product data"
// @Failure 404 {object} status_http.Response{data=string} "Product Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (p *productRoutes) getProductByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		p.handleResponse(c, status_http.BadRequest, "Product ID is required")
		return
	}

	product, err := p.productUscase.Get(c, id)
	if err != nil {
		p.handleResponse(c, status_http.InternalServerError, "error getting product")
		return
	}

	p.handleResponse(c, status_http.OK, product)
}

// @Router /product [put]
// @Summary Update an existing product
// @Description Update product details by ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param product body entity.UpdateProductRequest true "Product Details"
// @Failure 400 {object} status_http.Response{data=string} "Product updated successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (p *productRoutes) updateProduct(c *gin.Context) {
	var product entity.UpdateProductRequest

	if err := c.ShouldBindJSON(&product); err != nil {
		p.handleResponse(c, status_http.BadRequest, "invalid request")
		return
	}

	err := p.productUscase.Update(c, &product)
	if err != nil {
		p.handleResponse(c, status_http.InternalServerError, "error updating product")
		return
	}

	p.handleResponse(c, status_http.OK, "Product updated successfully")
}

// @Router /product/{id} [delete]
// @Summary Delete a product by ID
// @Description Delete a product from the database by ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} status_http.Response{data=string} "Product deleted successfully"
// @Success 400 {object} status_http.Response{data=string} "Bad request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (p *productRoutes) deleteProduct(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		p.handleResponse(c, status_http.BadRequest, "Product ID is required")
		return
	}

	err := p.productUscase.Delete(c, id)
	if err != nil {
		p.handleResponse(c, status_http.InternalServerError, "error deleting product")
		return
	}

	p.handleResponse(c, status_http.OK, "Product deleted successfully")
}

func (h *productRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
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
