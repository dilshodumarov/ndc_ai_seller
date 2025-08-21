package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/pkg/helper"
	"sugurta/internal/usecase/telegram"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InstagramRoutes struct {
	handlers.BaseHandler
	log      *zap.Logger
	cfg      *config.Config
	enforcer *casbin.CachedEnforcer
	telegram telegram.TelegramAccount
}

// NewAuthRoutes creates a new auth routes controller
func NewInstagramRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &InstagramRoutes{
		log:      option.Logger,
		cfg:      option.Config,
		enforcer: option.Enforcer,
		telegram: option.Telegram,
	}

	instagram := apiV1Group.Group("/instagram")
	{
		instagram.POST("/login", r.InstagramLogin)
		instagram.Any("/webhook/instagram", r.HandleInstagramWebhook)
		instagram.Any("/oauth/callback", r.HandleInstagramCallback)

	}
}

// InstagramLogin godoc
// @Summary Instagram login redirect
// @Description Redirect user to Instagram OAuth page
// @Security BearerAuth
// @Tags INSTAGRAM
// @Produce json
// @Success 302 {string} string "Redirect"
// @Router /instagram/login [post]
func (b *InstagramRoutes) InstagramLogin(c *gin.Context) {
	fmt.Println("Instagram login requested")

	businessID, code := helper.GetBusnessIdFromToken(c, b.cfg)
	if code != 0 {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	redirectURI := "https://dilshodforever.uz/v1/instagram/oauth/callback"

	authURL := fmt.Sprintf(
		"https://www.instagram.com/oauth/authorize?client_id=700909965624963&redirect_uri=%s&response_type=code&scope=instagram_business_basic,instagram_business_manage_messages,instagram_business_manage_comments,instagram_business_content_publish,instagram_business_manage_insights&state=%s",
		url.QueryEscape(redirectURI),
		url.QueryEscape(businessID),
	)

	c.JSON(200, gin.H{
		"url": authURL,
	})
}

func (b *InstagramRoutes) HandleInstagramWebhook(c *gin.Context) {
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
		fmt.Println("message: ", string(body))
		resp, err := http.Post("http://ai-seller-bot:8089/send-message-instagram", "application/json", bytes.NewReader(body))
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

func (b *InstagramRoutes) HandleInstagramCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		log.Println("Code not found in query params")
		c.JSON(http.StatusBadRequest, "Code not found in query params")
		return
	}
	state := c.Query("state")

	fmt.Println("Business ID:", state)
	// Step 1: Get short-lived access token
	data := url.Values{}
	// data.Set("client_id", b.cfg.AppConfig.ClientID)
	// data.Set("client_secret", b.cfg.AppConfig.ClientSecret)
	// data.Set("grant_type", b.cfg.AppConfig.GrantType)
	// data.Set("redirect_uri", b.cfg.AppConfig.RedirectURI)
	// data.Set("code", code)

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
		UserID      json.Number `json:"user_id"`
	}
	fmt.Println("Body: ", string(body))
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Println("Error parsing token response:", err)
		c.JSON(http.StatusInternalServerError, "Failed to parse token")
		return
	}

	// 🔥 Step 2: (Optional) Convert short-lived token to long-lived token
	longLivedToken, err := exchangeForLongLivedToken(tokenResp.AccessToken)
	if err != nil {
		log.Println("Error exchanging long-lived token:", err)
		c.JSON(http.StatusInternalServerError, "Failed to exchange long-lived token")
		return
	}

	// 🔥 Step 3: Subscribe to Webhook directly on Instagram
	subscribeURL := fmt.Sprintf("https://graph.instagram.com/v23.0/%s/subscribed_apps?subscribed_fields=comments,messages&access_token=%s", tokenResp.UserID.String(), longLivedToken)

	subscribeResp, err := http.Post(subscribeURL, "application/json", nil)
	if err != nil {
		log.Println("Error subscribing to webhook:", err)
		c.JSON(http.StatusInternalServerError, "Failed to subscribe to webhook")
		return
	}
	defer subscribeResp.Body.Close()

	subscribeBody, _ := io.ReadAll(subscribeResp.Body)
	log.Printf("Subscribed to Instagram Webhook: %s", string(subscribeBody))
	_, err = b.telegram.Create(c, entity.CreateTelegramAccountRequest{
		UserID:     tokenResp.UserID.String(),
		From:       "instagram",
		BusinessID: state,
	})
	if err != nil {

		log.Println("Error subscribing to webhook:", err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": longLivedToken,
		"user_id":      tokenResp.UserID,
		"message":      "Instagram connected and webhook subscribed",
	})
}

// Helper function to exchange long-lived token
func exchangeForLongLivedToken(shortToken string) (string, error) {
	url := fmt.Sprintf("https://graph.instagram.com/access_token?grant_type=ig_exchange_token&client_secret=%s&access_token=%s",
		"dff534402f4026921ee41af2f8a5c415", shortToken)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (h *InstagramRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
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
