// Package v1 implements routing paths. Each services in own file.
package api

import (
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	redisrepo "sugurta/internal/infrastructure/repository/redis"
	"sugurta/internal/pkg/config"
	botcomments "sugurta/internal/usecase/bot_comments"
	"sugurta/internal/usecase/business"
	"sugurta/internal/usecase/category"
	"sugurta/internal/usecase/chat"
	clienttype "sugurta/internal/usecase/client-type"
	"sugurta/internal/usecase/database"
	integration "sugurta/internal/usecase/intgration"
	"sugurta/internal/usecase/notification"
	"sugurta/internal/usecase/order"
	"sugurta/internal/usecase/product"
	"sugurta/internal/usecase/role"
	"sugurta/internal/usecase/settings"
	"sugurta/internal/usecase/telegram"
	"sugurta/internal/usecase/user"

	"github.com/casbin/casbin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sugurta/api/handlers"

	// Swagger docs.
	v1 "sugurta/api/handlers/v1"
	"sugurta/api/middleware"
	_ "sugurta/docs"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

type RouteOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Cache          redisrepo.Cache
	Enforcer       *casbin.CachedEnforcer

	User         user.User
	Role         role.Role
	ClientType   clienttype.ClientType
	Business     business.Business
	Product      product.Product
	Category     category.Category
	Order        order.Order
	Integration  integration.Integration
	BotComments  botcomments.BotCommandStorage
	Chat         chat.Chat
	Minio        *minio.Client
	Telegram     telegram.TelegramAccount
	Notification notification.Notification
	Settings     settings.SettingsStorage
	Database     database.Database
	// Service        grpcClients.ServiceClient
	// RefreshToken   refresh_token.RefreshToken
	// BrokerProducer event.BrokerProducer

}

// func NewRouter(app *gin.Engine, h Handler, auth usecase.Auth, product usecase.Product, integration usecase.Integration) {

// NewRouter -.
// Swagger spec:
// @title       AI Seller API
// @version     1.0
// @description This API powers the AI Seller platform for intelligent product recommendations and management.
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(option *RouteOption) *gin.Engine {
	handleOption := &handlers.HandlerOption{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Cache:          option.Cache,
		// Enforcer:       option.Enforcer,
		User:         option.User,
		Business:     option.Business,
		Order:        option.Order,
		Integration:  option.Integration,
		Product:      option.Product,
		Category:     option.Category,
		BotComments:  option.BotComments,
		Chat:         option.Chat,
		Minio:        option.Minio,
		Telegram:     option.Telegram,
		Notification: option.Notification,
		Settings:     option.Settings,
		Database:     option.Database,
	}

	app := gin.New()

	app.Use(gin.Logger(), gin.Recovery())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Authentication"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	app.Use(middleware.Logger(option.Logger))
	app.Use(middleware.Recovery(option.Logger))
	//e := casbin.NewEnforcer("internal/pkg/config/rbac.conf", "internal/pkg/config/policy.csv")
	e := casbin.NewEnforcer("/config/rbac.conf", "/config/policy.csv")
	app.Use(middleware.AuthMiddleware(e, *handleOption.Config))
	
	app.GET("/docs", func(ctx *gin.Context) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Web playground API",
			},
			Theme:    scalar.ThemeBluePlanet,
			DarkMode: true,
			Layout:   scalar.LayoutModern,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "error loading scalar docs",
			})
		}

		ctx.Header("Content-Type", "text/html")
		ctx.String(http.StatusOK, htmlContent)

	})

	// Swagger
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewAuthRoutes(apiV1Group, handleOption)
		v1.NewOrderRoutes(apiV1Group, handleOption)
		v1.NewBusinessRoutes(apiV1Group, handleOption)
		v1.NewIntegrationRoutes(apiV1Group, handleOption)
		// v1.NewClientTypeRoutes(apiV1Group, auth, h.cfg, h.log, h.inMemory)
		// v1.NewRoleRoutes(apiV1Group, auth, h.cfg, h.log, h.inMemory)
		v1.NewCategoryRoutes(apiV1Group, handleOption)
		v1.NewProductRoutes(apiV1Group, handleOption)
		v1.NewBotCommentsRoutes(apiV1Group, handleOption)
		v1.NewWebsocketRoutes(apiV1Group, handleOption)
		v1.NewTelegramRoutes(apiV1Group, handleOption)
		v1.NewMinioRoutes(apiV1Group, handleOption)
		v1.NewNotificationRoutes(apiV1Group, handleOption)
		v1.NewSettingsRoutes(apiV1Group, handleOption)
		v1.NewDatabaseRoutes(apiV1Group, handleOption)
		v1.NewInstagramRoutes(apiV1Group,handleOption)
	}

	return app
}
