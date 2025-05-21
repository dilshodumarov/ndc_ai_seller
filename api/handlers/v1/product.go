package v1

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/pkg/helper"
	"sugurta/internal/usecase/product"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
		productGroup.PUT("/update/:id", r.updateProduct)
		productGroup.PUT("/update", r.updateProducts)
		productGroup.GET("/list", r.ListProducts)
		productGroup.DELETE("/delete/:id", r.deleteProduct)
		productGroup.POST("/picture", r.addProductPicture)
		productGroup.GET("/export", r.ExportProducts)
		productGroup.GET("/export-csv", r.ExportProductsCSV)

	}

}

// @Router /product/create [post]
// @Summary Create a new product
// @Description Create a new product in the database
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body []entity.CreateProductRequestForSwagger true "Product Details"
// @Success 201 {object} status_http.Response{data=string} "Product created successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 401 {object} status_http.Response{data=string} "Unauthorized"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (p *productRoutes) createProduct(c *gin.Context) {
	var products []entity.CreateProductRequest

	if err := c.ShouldBindJSON(&products); err != nil {
		p.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	businessID, code := helper.GetBusnessIdFromToken(c, p.cfg)
	if code != 0 {
		p.handleResponse(c, status_http.Unauthorized, "Unauthorized")
		return
	}

	for _, product := range products {
		product.BusinessID = businessID
		if product.CategoryID != nil {
			if *product.CategoryID == "" {
				product.CategoryID = nil
			}
		}
		id, err := p.productUscase.Create(c, &product)
		if err != nil {
			p.handleResponse(c, status_http.InternalServerError, "Error while creating product")
			return
		}

		for _, imageURL := range product.Image_url {
			_, err := p.productUscase.AddPicture(c, &entity.CreateProductImage{
				ProductId: id,
				ImageUrl:  imageURL,
			})
			if err != nil {
				p.handleResponse(c, status_http.InternalServerError, "Error while saving product image")
				return
			}
		}
	}

	p.handleResponse(c, status_http.Created, "Product created successfully")

	// if product.Discount != 0 {
	// 	go func(product entity.CreateProductRequest, id string) {
	// 		botNotification := entity.BotNotification{
	// 			Guid:      product.BusinessID,
	// 			ProductId: id,
	// 		}
	// 		body, err := json.Marshal(botNotification)
	// 		if err != nil {
	// 			log.Println("Failed to marshal JSON:", err)
	// 			return
	// 		}

	// 		resp, err := http.Post("http://ai-seller-bot:8081/notification", "application/json", bytes.NewBuffer(body))
	// 		if err != nil {
	// 			log.Println("Failed to send request to bot:", err)
	// 			return
	// 		}
	// 		defer resp.Body.Close()

	// 		var BotResp entity.BotIntegrationResponse
	// 		if err := json.NewDecoder(resp.Body).Decode(&BotResp); err != nil {
	// 			log.Println("Failed to decode bot response:", err)
	// 			return
	// 		}
	// 		if BotResp.Code != 0 {
	// 			log.Println("Bot returned error:", BotResp.Message)
	// 		}
	// 	}(product, id)
	// }
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
		fmt.Println(err)
		p.handleResponse(c, status_http.InternalServerError, "error getting product")
		return
	}

	p.handleResponse(c, status_http.OK, product)
}

// @Router /product/list [get]
// @Summary List products
// @Description Get a paginated list of products with optional filters
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id query string false "Filter by Category ID"
// @Param product_id query int false "Filter by product id"
// @Param status query string false "Filter by status sotuvda or arxiv"
// @Param search query string false "Search in name, description, or short_info"
// @Param product_count query string false "filtr by product count"
// @Param limit query integer true "Number of products per page" default(10)
// @Param page query integer true "Page number " default(1)
// @Success 200 {object} status_http.Response{data=entity.GetAllProductsResponse} "List of Products"
// @Failure 400 {object} status_http.Response{data=string} "Bad request"
// @Failure 401 {object} status_http.Response{data=string} "Unauthorized"
// @Failure 500 {object} status_http.Response{data=string} "Internal server error"
func (p *productRoutes) ListProducts(c *gin.Context) {
	ownerId, code := helper.GetBusnessIdFromToken(c, p.cfg)
	if code != 0 {
		p.handleResponse(c, status_http.Unauthorized, "Unauthorized")
		return
	}

	filter := entity.ProductFilter{
		OwnerID:    ownerId,
		CategoryID: c.Query("category_id"),
		Search:     c.Query("search"),
		Status:     c.Query("status"),
	}
	prid := c.Query("product_id")
	if prid != "" {
		strprid, err := strconv.Atoi(prid)
		if err != nil {
			fmt.Println(err)
			p.handleResponse(c, status_http.BadRequest, "Invalid product id")
			return
		}
		filter.ProductId = strprid
	}
	count := c.Query("product_count")
	if count != "" {
		strcount, err := strconv.Atoi(count)
		if err != nil {
			fmt.Println(err)
			p.handleResponse(c, status_http.BadRequest, "Invalid profuct_count")
			return
		}
		filter.ProductCount = strcount
	}
	// Default values
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	fmt.Println(filter)
	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil || limit == 0 {
		p.handleResponse(c, status_http.BadRequest, "Invalid limit parameter")
		return
	}

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil || page == 0 {
		p.handleResponse(c, status_http.BadRequest, "Invalid page parameter")
		return
	}

	filter.Limit = limit
	filter.Page = page

	result, err := p.productUscase.List(c, filter)
	if err != nil {
		fmt.Println("List error:", err)
		p.handleResponse(c, status_http.InternalServerError, "Error listing products")
		return
	}

	p.handleResponse(c, status_http.OK, result)
}

// @Router /product/update/{id} [put]
// @Summary Update an existing product
// @Description Update product details by ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Security BearerAuth
// @Param product body entity.UpdateProductRequestForSwagger true "Product Details"
// @Failure 400 {object} status_http.Response{data=string} "Product updated successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (p *productRoutes) updateProduct(c *gin.Context) {
	var product entity.UpdateProductRequest
	id := c.Param("id")
	if err := c.ShouldBindJSON(&product); err != nil {
		p.handleResponse(c, status_http.BadRequest, "invalid request")
		return
	}

	product.ID = id
	err := p.productUscase.Update(c, &product)
	if err != nil {
		fmt.Println(err)
		p.handleResponse(c, status_http.InternalServerError, "error updating product")
		return
	}
	if len(product.Image_url) > 0 {
		err = p.productUscase.DeletePicture(c, id)
		if err != nil {
			fmt.Println(err)
			p.handleResponse(c, status_http.InternalServerError, "error while delete image")
			return
		}
		for i := 0; i < len(product.Image_url); i++ {
			_, err = p.productUscase.AddPicture(c, &entity.CreateProductImage{
				ProductId: id,
				ImageUrl:  product.Image_url[i],
			})
			if err != nil {
				fmt.Println(err)
				p.handleResponse(c, status_http.InternalServerError, "error while creating image")
				return
			}
		}
	}

	p.handleResponse(c, status_http.OK, "Product updated successfully")
}

// @Router /product/update [put]
// @Summary Update multiple products
// @Description Update details of multiple products
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param products body []entity.UpdateProductRequest true "List of products to update"
// @Success 200 {object} status_http.Response{data=string} "Products updated successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (p *productRoutes) updateProducts(c *gin.Context) {
	var products []entity.UpdateProductRequest

	if err := c.ShouldBindJSON(&products); err != nil {
		p.handleResponse(c, status_http.BadRequest, "invalid request")
		return
	}

	for _, product := range products {
		if product.ID == "" {
			p.handleResponse(c, status_http.BadRequest, "product ID is required")
			return
		}

		err := p.productUscase.Update(c, &product)
		if err != nil {
			p.handleResponse(c, status_http.InternalServerError, fmt.Sprintf("error updating product with ID %s", product.ID))
			return
		}

		if len(product.Image_url) > 0 {
			// Delete old images
			err = p.productUscase.DeletePicture(c, product.ID)
			if err != nil {
				p.handleResponse(c, status_http.InternalServerError, fmt.Sprintf("error deleting images for product %s", product.ID))
				return
			}

			// Add new images
			for _, imageURL := range product.Image_url {
				_, err = p.productUscase.AddPicture(c, &entity.CreateProductImage{
					ProductId: product.ID,
					ImageUrl:  imageURL,
				})
				if err != nil {
					p.handleResponse(c, status_http.InternalServerError, fmt.Sprintf("error adding image to product %s", product.ID))
					return
				}
			}
		}
	}

	p.handleResponse(c, status_http.OK, "Products updated successfully")
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
		fmt.Println(err)
		p.handleResponse(c, status_http.InternalServerError, "error deleting product")
		return
	}

	p.handleResponse(c, status_http.OK, "Product deleted successfully")
}

// @Router /product/picture [post]
// @Summary Add a picture to a product
// @Description Attach an image to an existing product by product ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param picture body entity.CreateProductImage true "Product Picture Details"
// @Success 201 {object} status_http.Response{data=string} "Image added successfully"
// @Failure 400 {object} status_http.Response{data=string} "Invalid request data"
// @Failure 500 {object} status_http.Response{data=string} "Failed to add image"
func (p *productRoutes) addProductPicture(c *gin.Context) {
	var req entity.CreateProductImage

	if err := c.ShouldBindJSON(&req); err != nil {
		p.handleResponse(c, status_http.BadRequest, "Invalid request: "+err.Error())
		return
	}

	id, err := p.productUscase.AddPicture(c, &req)
	if err != nil {
		p.handleResponse(c, status_http.InternalServerError, "Failed to add image: "+err.Error())
		return
	}

	p.handleResponse(c, status_http.Created, id)
}

// @Router /product/export [get]
// @Summary Export Products to Excel
// @Description Export product list with filters in Excel format
// @Tags PRODUCT
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Param category_id query string false "Filter by Category ID"
// @Param product_id query int false "Filter by product id"
// @Param status query string false "Filter by status sotuvda or arxiv"
// @Param search query string false "Search in name, description, or short_info"
// @Param product_count query string false "Filter by product count"
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} status_http.Response{data=string}
// @Failure 401 {object} status_http.Response{data=string}
// @Failure 500 {object} status_http.Response{data=string}
func (p *productRoutes) ExportProducts(c *gin.Context) {
	ownerId, code := helper.GetBusnessIdFromToken(c, p.cfg)
	if code != 0 {
		p.handleResponse(c, status_http.Unauthorized, "Unauthorized")
		return
	}

	filter := entity.ProductFilter{
		OwnerID:    ownerId,
		CategoryID: c.Query("category_id"),
		Search:     c.Query("search"),
		Status:     c.Query("status"),
	}
	prid := c.Query("product_id")
	if prid != "" {
		id, err := strconv.Atoi(prid)
		if err != nil {
			p.handleResponse(c, status_http.BadRequest, "Invalid product id")
			return
		}
		filter.ProductId = id
	}
	count := c.Query("product_count")
	if count != "" {
		cInt, err := strconv.Atoi(count)
		if err != nil {
			p.handleResponse(c, status_http.BadRequest, "Invalid product count")
			return
		}
		filter.ProductCount = cInt
	}

	// Excelga eksport qilish
	fileBytes, err := p.exportToExcel(c, &filter)
	if err != nil {
		p.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename=products.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// @Router /product/export-csv [get]
// @Summary Export Products to CSV
// @Description Export product list with filters in CSV format
// @Tags PRODUCT
// @Accept json
// @Produce text/csv
// @Security BearerAuth
// @Param category_id query string false "Filter by Category ID"
// @Param product_id query int false "Filter by product id"
// @Param status query string false "Filter by status sotuvda or arxiv"
// @Param search query string false "Search in name, description, or short_info"
// @Param product_count query string false "Filter by product count"
// @Success 200 {file} file "CSV file"
// @Failure 400 {object} status_http.Response{data=string}
// @Failure 401 {object} status_http.Response{data=string}
// @Failure 500 {object} status_http.Response{data=string}
func (p *productRoutes) ExportProductsCSV(c *gin.Context) {
	ownerId, code := helper.GetBusnessIdFromToken(c, p.cfg)
	if code != 0 {
		p.handleResponse(c, status_http.Unauthorized, "Unauthorized")
		return
	}

	filter := entity.ProductFilter{
		OwnerID:    ownerId,
		CategoryID: c.Query("category_id"),
		Search:     c.Query("search"),
		Status:     c.Query("status"),
	}
	prid := c.Query("product_id")
	if prid != "" {
		id, err := strconv.Atoi(prid)
		if err != nil {
			p.handleResponse(c, status_http.BadRequest, "Invalid product id")
			return
		}
		filter.ProductId = id
	}
	count := c.Query("product_count")
	if count != "" {
		cInt, err := strconv.Atoi(count)
		if err != nil {
			p.handleResponse(c, status_http.BadRequest, "Invalid product count")
			return
		}
		filter.ProductCount = cInt
	}

	// CSVga eksport qilish
	fileBytes, err := p.exportToCSV(c, &filter)
	if err != nil {
		p.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename=products.csv")
	c.Data(http.StatusOK, "text/csv", fileBytes)
}

func (p *productRoutes) exportToCSV(ctx context.Context, filter *entity.ProductFilter) ([]byte, error) {
	filter.Limit = 0
	filter.Page = 0

	result, err := p.productUscase.List(ctx, *filter)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header
	headers := []string{
		"Product ID", "Business ID", "Name", "Category ID", "Category Name",
		"Short Info", "Description", "Cost", "Count", "Discount", "Discount Cost",
		"Status", "Created At", "Updated At",
	}
	writer.Write(headers)

	// Data rows
	for _, prod := range result.Items {
		row := []string{
			prod.ID,
			prod.BusinessID,
			prod.Name,
			prod.CategoryID,
			prod.CategoryName,
			prod.ShortInfo,
			prod.Description,
			fmt.Sprintf("%.d", prod.Cost),
			strconv.Itoa(prod.Count),
			fmt.Sprintf("%.d", prod.Discount),
			fmt.Sprintf("%.d", prod.DiscountCost),
			func() string {
				if prod.Status {
					return "Sotuvda"
				}
				return "Arxiv"
			}(),
			prod.CreatedAt.Format("2006-01-02 15:04:05"),
			prod.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *productRoutes) exportToExcel(ctx context.Context, filter *entity.ProductFilter) ([]byte, error) {
	// limit=0, page=0 → barcha mahsulotlarni olish
	filter.Limit = 0
	filter.Page = 0

	result, err := p.productUscase.List(ctx, *filter)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Products"
	f.SetSheetName("Sheet1", sheet)

	// Header row
	headers := []string{
		"Product ID", "Business ID", "Name", "Category ID", "Category Name",
		"Short Info", "Description", "Cost", "Count", "Discount", "Discount Cost",
		"Status", "Created At", "Updated At",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Data rows
	for idx, prod := range result.Items {
		values := []interface{}{
			prod.ID,
			prod.BusinessID,
			prod.Name,
			prod.CategoryID,
			prod.CategoryName,
			prod.ShortInfo,
			prod.Description,
			prod.Cost,
			prod.Count,
			prod.Discount,
			prod.DiscountCost,
			func() string {
				if prod.Status {
					return "Sotuvda"
				}
				return "Arxiv"
			}(),
			prod.CreatedAt.Format("2006-01-02 15:04:05"),
			prod.UpdatedAt.Format("2006-01-02 15:04:05"),
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

// if product.Discount != 0 {
// 	BusinessID, code := helper.GetBusnessIdFromToken(c, p.Config)
// 	if code != 0 {
// 		p.handleResponse(c, status_http.Unauthorized, "Unauthorized")
// 	}
// 	botNotification := entity.BotNotification{
// 		Guid:      BusinessID,
// 		ProductId: product.ID,
// 	}
// 	body, err := json.Marshal(botNotification)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
// 		return
// 	}
// 	resp, err := http.Post("http://ai-seller-bot:8081/notification", "application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
// 		return
// 	}
// 	var BotResp entity.BotIntegrationResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&BotResp); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
// 		return
// 	}
// 	if BotResp.Code != 0 {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bot start" + BotResp.Message})
// 		return
// 	}
// 	p.handleResponse(c, status_http.OK, "Product updated successfully")
// 	return
// }
