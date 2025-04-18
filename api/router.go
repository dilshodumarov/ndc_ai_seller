// Package v1 implements routing paths. Each services in own file.
package api

import (
	"net/http"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	redisrepo "sugurta/internal/infrastructure/repository/redis"
	"sugurta/internal/pkg/config"
	botcomments "sugurta/internal/usecase/bot_comments"
	"sugurta/internal/usecase/business"
	"sugurta/internal/usecase/category"
	clienttype "sugurta/internal/usecase/client-type"
	integration "sugurta/internal/usecase/intgration"
	"sugurta/internal/usecase/order"
	"sugurta/internal/usecase/product"
	"sugurta/internal/usecase/role"
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

	User        user.User
	Role        role.Role
	ClientType  clienttype.ClientType
	Business    business.Business
	Product     product.Product
	Category    category.Category
	Order       order.Order
	Integration integration.Integration
	BotComments botcomments.BotCommandStorage
	// Service        grpcClients.ServiceClient
	// RefreshToken   refresh_token.RefreshToken
	// BrokerProducer event.BrokerProducer

}

// func NewRouter(app *gin.Engine, h Handler, auth usecase.Auth, product usecase.Product, integration usecase.Integration) {

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
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
		User:           option.User,
		Business:       option.Business,
		Order:          option.Order,
		Integration:    option.Integration,
		Product:        option.Product,
		Category:       option.Category,
		BotComments:    option.BotComments,
	}

	app := gin.New()

	app.Use(gin.Logger(), gin.Recovery())

	// Options
	app.Use(middleware.Logger(option.Logger))
	app.Use(middleware.Recovery(option.Logger))
	e := casbin.NewEnforcer("internal/pkg/config/rbac.conf", "internal/pkg/config/policy.csv")
	app.Use(middleware.AuthMiddleware(e,*handleOption.Config))
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-Id"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
	    v1.NewCategoryRoutes(apiV1Group,handleOption)
		v1.NewProductRoutes(apiV1Group, handleOption)
		v1.NewBotCommentsRoutes(apiV1Group,handleOption)
		v1.NewWebsocketRoutes(apiV1Group,handleOption)
	}

	return app
}
