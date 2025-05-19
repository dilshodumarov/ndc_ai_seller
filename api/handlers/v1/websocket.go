package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sugurta/api/handlers"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sugurta/internal/usecase/chat"
	"sugurta/internal/usecase/notification"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type websocketRoutes struct {
	handlers.BaseHandler
	log              *zap.Logger
	cfg              *config.Config
	enforcer         *casbin.CachedEnforcer
	clients          map[string]*websocket.Conn
	mu               sync.Mutex
	ChatRepo         chat.Chat
	NotificationRepo notification.Notification
}

func NewWebsocketRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &websocketRoutes{
		log:              option.Logger,
		cfg:              option.Config,
		enforcer:         option.Enforcer,
		ChatRepo:         option.Chat,
		clients:          make(map[string]*websocket.Conn),
		NotificationRepo: option.Notification,
	}

	websocketGroup := apiV1Group.Group("/websocket")
	{
		websocketGroup.GET("chat/connect", r.WebSocketHandler)
		websocketGroup.POST("/chat/send-message", r.SendChatMessage)
	}
	apiV1Group.GET("/chat/list/:bussnesid", r.ListChat)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *websocketRoutes) WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("WebSocket upgrade error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade failed"})
		return
	}

	userID := c.Query("id")
	if userID == "" {
		h.log.Error("User id is nill")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		conn.Close()
		return
	}

	h.AddClient(userID, conn)
	h.log.Info("New WebSocket client connected", zap.String("userID", userID))

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				h.log.Error("Error reading message", zap.Error(err))
				h.mu.Lock()
				delete(h.clients, userID)
				h.mu.Unlock()
				conn.Close()
				break
			}
			resp, err := http.Post("http://ai-seller-bot:8081/send-message", "application/json", bytes.NewBuffer(msg))
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
			h.log.Info("Message send succesfully", zap.String("message", string(msg)))

			// optional: DB ga saqlash, boshqa foydalanuvchiga yuborish, logger, AI javob, va hokazo

		}
	}()
}

func (h *websocketRoutes) AddClient(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = conn
}

func (h *websocketRoutes) SendChatMessage(c *gin.Context) {
	var msg entity.SendMessageResponse
	var UserID string
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	if msg.Type == "chat" {
		UserID = msg.ChatMessage.UserId
	} else {
		UserID = msg.Notifications.UserId
		_,err:=h.NotificationRepo.Create(c, &entity.CreateNotificationRequest{
			UserID:  msg.Notifications.UserId,
			Title:   msg.Notifications.Title,
			Message: msg.Notifications.Content,
			Type:    "info",
		})
		if err!=nil{
			fmt.Println(err)
		}
	}
	h.mu.Lock()
	conn, exists := h.clients[UserID]
	h.mu.Unlock()

	if !exists {
		h.log.Warn("No WebSocket connection for user", zap.String("userid", UserID))
		c.JSON(http.StatusNotFound, gin.H{"error": "No WebSocket connection for this user"})
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		h.log.Error("Failed to marshal chat message", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process message"})
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		h.log.Error("WebSocket send error", zap.Error(err))

		h.mu.Lock()
		delete(h.clients, UserID)
		h.mu.Unlock()

		conn.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket send error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

// ListChatHistory godoc
// @Summary Get list of chat history by ChatID and BusinessID
// @Description Retrieve chat messages (user or AI) by a given Chat ID and Business ID
// @Tags CHAT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chatid query int false "Chat ID"
// @Param bussnesid path string true "Business ID"
// @Param limit query int false "Maximum number of chat messages to return. If 0 or not provided, all messages will be returned. Default: 100"
// @Success 200 {object} status_http.Response{data=[]entity.SendMessage} "Chat history list"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @Router /chat/list/{bussnesid} [get]
func (h *websocketRoutes) ListChat(c *gin.Context) {
	chatid := c.Query("chatid")
	bussnesid := c.Param("bussnesid")
	limitStr := c.DefaultQuery("limit", "0")
	limit, _ := strconv.Atoi(limitStr)
	var res entity.ListChatHistoryRequest
	if bussnesid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "businessid are required"})
		return
	}
	if chatid != "" {
		intchatid, err := strconv.ParseInt(chatid, 10, 64)
		res.ChatID = intchatid
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chatid"})
			return
		}
	}
	res.BusinessID = bussnesid
	res.Limit = limit

	chats, err := h.ChatRepo.List(c, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}
