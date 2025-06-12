package v1

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/pkg/helper"
	"sugurta/internal/usecase/business"
	"sugurta/internal/usecase/settings"

	"net/http"
	"net/url"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type businessRoutes struct {
	handlers.BaseHandler
	log            *zap.Logger
	cfg            *config.Config
	enforcer       *casbin.CachedEnforcer
	bussnesUscase  business.Business
	settingsUscase settings.SettingsStorage
}

// NewAuthRoutes creates a new auth routes controller
func NewBusinessRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &businessRoutes{
		log:            option.Logger,
		cfg:            option.Config,
		enforcer:       option.Enforcer,
		bussnesUscase:  option.Business,
		settingsUscase: option.Settings,
	}

	business := apiV1Group.Group("/business")
	{
		business.POST("/create", r.CreateBusiness)
		business.GET("/get/:id", r.GetBusiness)
		business.PUT("/update/:id", r.UpdateBusiness)
		business.DELETE("/delete/:id", r.DeleteBusiness)
		business.GET("/list", r.GetAllBusinesses)
		business.Any("/webhook/instagram", r.HandleInstagramWebhook)
		business.Any("/oauth/callback", r.HandleInstagramCallback)

	}
}

// CreateBusiness godoc
// @Summary Create a new business
// @Description Create a new business with the provided details
// @Tags BUSINESS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param business body entity.CreateBusinessRequest true "Business details"
// @Success 201 {object} status_http.Response{data=string} "Success"
// @Response 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /business/create [post]
func (b *businessRoutes) CreateBusiness(c *gin.Context) {
	var business entity.CreateBusinessRequest

	if err := c.ShouldBindJSON(&business); err != nil {
		b.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}
	UserId, code := helper.GetUserIdFromToken(c, b.cfg)
	if code != 0 {
		b.handleResponse(c, status_http.Unauthorized, "Unauthorized")
	}
	business.OwnerID = UserId
	if business.Name == "" {
		b.handleResponse(c, status_http.BadRequest, "name is required")
	}
	if business.OwnerID == "" {
		b.handleResponse(c, status_http.BadRequest, "description is required")
	}

	id, err := b.bussnesUscase.Create(c, &business)
	if err != nil {
		b.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	err = b.settingsUscase.CreateDefaultOrderStatuses(c, id)
	if err != nil {
		fmt.Println(err)
		b.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	err = b.settingsUscase.CreateSettings(c, &entity.CreateSettingsRequest{BusinessID: id})
	if err != nil {
		fmt.Println(err)
		b.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	b.handleResponse(c, status_http.Created, "Business created successfully")
}

// GetBusiness godoc
// @Summary Get a business by ID
// @Description Get business details by its ID
// @Tags BUSINESS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID"
// @Success 200 {object} status_http.Response{data=string} "Business"
// @Response400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /business/get/{id} [get]
func (b *businessRoutes) GetBusiness(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		b.handleResponse(c, status_http.BadRequest, "business ID is required")
		return
	}

	business, err := b.bussnesUscase.Get(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			b.handleResponse(c, status_http.BadRequest, "business not found")
			return
		}
		b.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	b.handleResponse(c, status_http.OK, business)
}

// UpdateBusiness godoc
// @Summary Update a business
// @Description Update business details
// @Tags BUSINESS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID"
// @Param business body entity.UpdateBusinessRequestForSwagger true "Business details"
// @Success 200 {object} status_http.Response{data=string} "Business updated"
// @Response 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /business/update/{id} [put]
func (b *businessRoutes) UpdateBusiness(c *gin.Context) {
	var business entity.UpdateBusinessRequest
	if err := c.ShouldBindJSON(&business); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	business.ID = id
	if err := b.bussnesUscase.Update(c, &business); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, business)
}

// DeleteBusiness godoc
// @Summary Delete a business
// @Description Delete a business by ID
// @Tags BUSINESS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID"
// @Success 200 {object} status_http.Response{data=string} "Successfully deleted"
// @Response 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /business/delete/{id} [delete]
func (b *businessRoutes) DeleteBusiness(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		b.handleResponse(c, status_http.BadRequest, "business ID is required")
		return
	}

	if err := b.bussnesUscase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Business successfully deleted"})
}

// GetAllBusinesses godoc
// @Summary Get all businesses
// @Description Get all businesses with pagination
// @Tags BUSINESS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit per page" default(10)
// @Param page query int false "Page number" default(0)
// @Param user_id query string false "UserID"
// @Success 200 {object} status_http.Response
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /business/list [get]
func (b *businessRoutes) GetAllBusinesses(c *gin.Context) {
	limitStr := c.Query("limit")
	pageStr := c.Query("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 // default value
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1 // default value
	}

	req := entity.GetAllBusinessesRequest{
		Limit:   limit,
		Page:    page,
		OwnerID: c.Query("user_id"),
	}

	fmt.Println(req)
	// UserId, code := helper.GetUserIdFromToken(c, b.Config)
	// if code != 0 {

	// 	b.handleResponse(c, status_http.Unauthorized, "Unauthorized")
	// }
	// req.OwnerID=UserId
	businesses, err := b.bussnesUscase.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, businesses)
}

func (h *businessRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
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

func (b *businessRoutes) HandleInstagramWebhook(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		challenge := c.Query("hub.challenge")
		verificationToken := c.Query("hub.verify_token")

		if verificationToken == "your_verification_token" {
			c.String(http.StatusOK, challenge)
		} else {
			b.handleResponse(c, status_http.Status{
				Code:          http.StatusForbidden,
				Status:        "error",
				Description:   "Invalid verification token",
				CustomMessage: "Instagram verification failed",
			}, nil)
		}

	case http.MethodPost:
		defer c.Request.Body.Close()
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			b.handleResponse(c, status_http.Status{
				Code:          http.StatusInternalServerError,
				Status:        "error",
				Description:   "Error reading request body",
				CustomMessage: err.Error(),
			}, nil)
			return
		}
		fmt.Println("message: ",string(body))
		resp, err := http.Post("http://ai-seller-bot:8081/send-message-instagram", "application/json", bytes.NewReader(body))
		if err != nil {
			b.handleResponse(c, status_http.Status{
				Code:          http.StatusInternalServerError,
				Status:        "error",
				Description:   "Failed to send request to bot",
				CustomMessage: err.Error(),
			}, nil)
			return
		}
		defer resp.Body.Close()

		b.handleResponse(c, status_http.Status{
			Code:        http.StatusOK,
			Status:      "success",
			Description: "Message forwarded to bot successfully",
		}, nil)

	default:
		b.handleResponse(c, status_http.Status{
			Code:        http.StatusMethodNotAllowed,
			Status:      "error",
			Description: "Method not allowed",
		}, nil)
	}
}

func (b *businessRoutes) HandleInstagramCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		log.Println("Code not found in query params")
		c.JSON(http.StatusBadRequest, "Code not found in query params")
		return
	}

	data := url.Values{}
	data.Set("client_id", "700909965624963")
	data.Set("client_secret", "22f5cd4d15549e880f749b1332b503fd")
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "https://dilshodforever.uz/v1/business/oauth/callback")
	data.Set("code", code)

	resp, err := http.PostForm("https://api.instagram.com/oauth/access_token", data)
	if err != nil {
		log.Println("Error requesting access token:", err)
		c.JSON(http.StatusInternalServerError, "Failed to request access token")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var tokenResp struct {
		AccessToken string      `json:"access_token"`
		UserID      json.Number `json:"user_id"` // json.Number bu string/raqam formatda o'qiy oladi
	}
	fmt.Println("Body: ", string(body))
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Println("Error parsing token response:", err)
		c.JSON(http.StatusInternalServerError, "Failed to parse token")
		return
	}

	// 🔥 Step 2: Get Facebook pages linked to the user
	pagesURL := fmt.Sprintf("https://graph.facebook.com/v19.0/me/accounts?access_token=%s", tokenResp.AccessToken)
	pagesResp, err := http.Get(pagesURL)
	if err != nil {
		log.Println("Error getting pages:", err)
		c.JSON(http.StatusInternalServerError, "Failed to get pages")
		return
	}
	defer pagesResp.Body.Close()

	var pages struct {
		Data []struct {
			ID                     string `json:"id"`
			Name                   string `json:"name"`
			AccessToken            string `json:"access_token"`
			ConnectedInstagramAcct struct {
				ID string `json:"id"`
			} `json:"connected_instagram_account"`
		} `json:"data"`
	}

	pagesBody, _ := io.ReadAll(pagesResp.Body)
	if err := json.Unmarshal(pagesBody, &pages); err != nil {
		log.Println("Error parsing pages:", err)
		c.JSON(http.StatusInternalServerError, "Failed to parse pages")
		return
	}
	id := ""
	// 🔥 Step 3: Subscribe each page to webhook
	for _, page := range pages.Data {
		if page.ConnectedInstagramAcct.ID == "" {
			log.Printf("Page '%s' has no connected Instagram account", page.Name)
			continue
		}
		fmt.Println(1111111, page.ID)
		id = page.ID
		subscribeURL := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/subscribed_apps", page.ID)
		subscribeResp, err := http.PostForm(subscribeURL, url.Values{
			"access_token": {page.AccessToken},
		})
		if err != nil {
			log.Printf("Error subscribing page %s: %v", page.Name, err)
			continue
		}
		defer subscribeResp.Body.Close()

		subBody, _ := io.ReadAll(subscribeResp.Body)
		log.Printf("Subscribed to page '%s': %s", page.Name, string(subBody))
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenResp.AccessToken,
		"user_id":      tokenResp.UserID,
		"message":      "Instagram connected and webhook subscribed",
		"pageid":       id,
	})
}



