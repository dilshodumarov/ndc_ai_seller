package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sugurta/api/handlers"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type websocketRoutes struct {
	handlers.BaseHandler
	log      *zap.Logger
	cfg      *config.Config
	enforcer *casbin.CachedEnforcer
	clients  map[string]*websocket.Conn
	mu       sync.Mutex
}

func NewWebsocketRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &websocketRoutes{
		log:      option.Logger,
		cfg:      option.Config,
		enforcer: option.Enforcer,
		clients:  make(map[string]*websocket.Conn),
	}

	websocketGroup := apiV1Group.Group("/websocket")
	{
		websocketGroup.GET("chat/connect", r.WebSocketHandler)
		websocketGroup.POST("/chat/send-message", r.SendChatMessage)
	}
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
			//h.log.Info("Received message from frontend", zap.String("message", chat.Message))

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
	var msg entity.ChatHistory
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	h.mu.Lock()
	conn, exists := h.clients[msg.BusinessId]
	h.mu.Unlock()

	if !exists {
		h.log.Warn("No WebSocket connection for user", zap.String("platformID", msg.PlatformID))
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
		delete(h.clients, msg.PlatformID)
		h.mu.Unlock()

		conn.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket send error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}
