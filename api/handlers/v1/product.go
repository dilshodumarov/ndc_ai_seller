package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
		productGroup.POST("/create", r.createProduct)
		productGroup.GET("/get/:id", r.getProductByID)
		productGroup.PUT("/update", r.updateProduct)
		productGroup.DELETE("/delete/:id", r.deleteProduct)
	}

}

// @Router /product/create [post]
// @Summary Create a new product
// @Description Create a new product in the database
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body entity.CreateProductRequest true "Product Details"
// @Success 201 {object} status_http.Response{data=string} "Product created successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (p *productRoutes) createProduct(c *gin.Context) {
	var product entity.CreateProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		p.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	id, err := p.productUscase.Create(c, &product)
	if err != nil {
		p.handleResponse(c, status_http.InternalServerError, "error while creating product")
		return
	}

	if product.Discount != 0 {
		botNotification := entity.BotNotification{
			Guid:      product.BusinessID,
			ProductId: id,
		}
		body, err := json.Marshal(botNotification)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
			return
		}
		resp, err := http.Post("http://localhost:8081/notification", "application/json", bytes.NewBuffer(body))
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
	}
	p.handleResponse(c, status_http.Created, "Product created successfully")
}

// @Router /product/get/{id} [get]
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Security BearerAuth
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

// @Router /product/update [put]
// @Summary Update an existing product
// @Description Update product details by ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Security BearerAuth
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
		fmt.Println(err)
		p.handleResponse(c, status_http.InternalServerError, "error updating product")
		return
	}
	if product.Discount != 0 {
		botNotification := entity.BotNotification{
			Guid:      product.BusinessID,
			ProductId: product.ID,
		}
		body, err := json.Marshal(botNotification)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
			return
		}
		resp, err := http.Post("http://localhost:8081/notification", "application/json", bytes.NewBuffer(body))
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
		p.handleResponse(c, status_http.OK, "Product updated successfully")
		return
	}
	p.handleResponse(c, status_http.OK, "Product updated successfully")
}

// @Router /product/delete/{id} [delete]
// @Summary Delete a product by ID
// @Description Delete a product from the database by ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Security BearerAuth
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
