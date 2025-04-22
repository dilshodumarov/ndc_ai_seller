// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	"net/http"
	"sugurta/internal/infrastructure/repository/postgresql"
	"sugurta/internal/pkg/config"
	"sugurta/internal/pkg/policy"
	"sugurta/internal/pkg/redis"
	botcomments "sugurta/internal/usecase/bot_comments"
	"sugurta/internal/usecase/business"
	"sugurta/internal/usecase/category"
	"sugurta/internal/usecase/chat"
	clienttype "sugurta/internal/usecase/client-type"
	integration "sugurta/internal/usecase/intgration"
	"sugurta/internal/usecase/order"
	"sugurta/internal/usecase/product"
	"sugurta/internal/usecase/role"
	"sugurta/internal/usecase/user"

	"sugurta/api"

	redisrepo "sugurta/internal/infrastructure/repository/redis"
	"sugurta/internal/pkg/logger"
	"sugurta/internal/pkg/postgres"
)

type App struct {
	Config   *config.Config
	Logger   *zap.Logger
	DB       *postgres.Postgres
	RedisDB  *redis.RedisDB
	server   *http.Server
	Enforcer *casbin.CachedEnforcer
	// Clients        grpcService.ServiceClient
	ShutdownOTLP func() error
	// BrokerProducer event.BrokerProducer
	user        user.User
	role        role.Role
	clientType  clienttype.ClientType
	business    business.Business
	product     product.Product
	category    category.Category
	order       order.Order
	integration integration.Integration
	BotComments botcomments.BotCommandStorage
	Chat        chat.Chat
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.Log.Level, cfg.Environment, cfg.App.Name+".log")
	if err != nil {
		return nil, err
	}

	// postgres init
	db, err := postgres.New(&cfg, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		return nil, err
	}

	// redis init
	redisdb, err := redis.New(&cfg)
	if err != nil {
		return nil, err
	}

	// otlp collector init
	// shutdownOTLP, err := otlp.InitOTLPProvider(&cfg)
	// if err != nil {
	// 	return nil, err
	// }

	// initialization enforcer
	enforcer, err := policy.NewCachedEnforcer(&cfg, logger)
	if err != nil {
		return nil, err
	}

	fmt.Println("here 1")

	enforcer.SetCache(policy.NewCache(&redisdb.Client))

	fmt.Println("here 2")

	var (
		contextTimeout time.Duration
	)

	// context timeout initialization
	contextTimeout, err = time.ParseDuration(cfg.Context.Timeout)
	if err != nil {
		return nil, err
	}

	userRepo := postgresql.NewUserRepo(db)
	userUseCase := user.NewUserService(contextTimeout, userRepo)
	fmt.Println("user repo: ", userUseCase)
	roleRepo := postgresql.NewRoleRepo(db)
	roleUseCase := role.NewRoleService(contextTimeout, roleRepo)

	businessRepo := postgresql.NewBusinessRepo(db)
	businessUseCase := business.NewbusinessService(contextTimeout, businessRepo)

	clientTypeRepo := postgresql.NewClientTypeRepo(db)
	clientTypeUseCase := clienttype.NewClientTypeService(contextTimeout, clientTypeRepo)

	productRepo := postgresql.NewProductRepo(db)
	productUseCase := product.NewProductService(contextTimeout, productRepo)

	categoryRepo := postgresql.NewCategoryRepo(db)
	categoryUseCase := category.NewCategoryService(contextTimeout, categoryRepo)

	intgrationRepo := postgresql.NewIntegrationRepo(db)
	intgrationUscase := integration.NewIntegrationService(contextTimeout, intgrationRepo)

	orderRepo := postgresql.NewOrderRepo(db)
	orderUseCase := order.NewRoleService(contextTimeout, orderRepo)

	botCommentsRepo := postgresql.NewBotCommandsRepo(db)
	botCommentsUscase := botcomments.NewbotCommentsService(contextTimeout, botCommentsRepo)

	ChatRepo := postgresql.NewChatRepo(db)
	ChatUscase := chat.NewChatService(contextTimeout, ChatRepo)

	fmt.Println("order repo: ", orderUseCase)
	fmt.Println("here 3")

	return &App{
		Config:   &cfg,
		Logger:   logger,
		DB:       db,
		RedisDB:  redisdb,
		Enforcer: enforcer,
		//ShutdownOTLP: shutdownOTLP,
		user:        userUseCase,
		role:        roleUseCase,
		business:    businessUseCase,
		clientType:  clientTypeUseCase,
		product:     productUseCase,
		category:    categoryUseCase,
		order:       orderUseCase,
		integration: intgrationUscase,
		BotComments: botCommentsUscase,
		Chat:         ChatUscase,
		// BrokerProducer: event.NewBrokerProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic),
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	fmt.Println("here 4")
	// clients, err := grpcService.New(a.Config)
	// if err != nil {
	// 	return err
	// }
	// a.Clients = clients

	// initialize cache
	cache := redisrepo.NewCache(a.RedisDB)
	fmt.Println("here 5")

	// tokenRepo := postgresql.NewRefreshTokenRepo(a.DB)

	// initialize token service
	// refreshTokenService := refresh_token.NewRefreshTokenService(contextTimeout, tokenRepo)

	// api init
	handler := api.NewRouter(&api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		Cache:          cache,
		// Enforcer:       a.Enforcer,
		User:        a.user,
		Business:    a.business,
		Order:       a.order,
		Integration: a.integration,
		Product:     a.product,
		Category:    a.category,
		BotComments: a.BotComments,
		Chat:        a.Chat,
	})

	fmt.Println("here 5")

	if err = a.Enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("error during enforcer load policy: %w", err)
	}

	fmt.Println("here 5")

	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close database
	a.DB.Close()

	// close grpc connections
	// a.Clients.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// shutdown otlp collector
	if err := a.ShutdownOTLP(); err != nil {
		a.Logger.Error("shutdown otlp collector", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}

// Run creates objects via constructors.
// func Run(cfg *config.Config) {
// 	// Logger
// 	l := logger.New(cfg.Log.Level)

// 	// Repository
// 	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
// 	if err != nil {
// 		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
// 	}
// 	defer pg.Close()

// 	authUsecase := auth.NewAuthUseCase(
// 		persistent.NewAuthRepo(pg),
// 		// persistent.NewProductRepo(pg),
// 		// persistent.NewIntegrationRepo(pg),
// 		l,
// 	)

// 	businessUsecase := business.NewBusinessUseCase(
// 		persistent.NewBusinessRepo(pg),
// 		l,
// 	)

// 	// productUseCase := product.NewProductUseCase(
// 	// 	persistent.NewAuthRepo(pg),
// 	// 	persistent.NewProductRepo(pg),
// 	// 	persistent.NewIntegrationRepo(pg),
// 	// 	l,
// 	// )

// 	// integrationUseCase := integration.NewIntegrationUseCase(
// 	// 	persistent.NewAuthRepo(pg),
// 	// 	persistent.NewProductRepo(pg),
// 	// 	persistent.NewIntegrationRepo(pg),
// 	// 	l,
// 	// )

// 	cache := cache.NewInMemoryStorage(cfg)

// 	// HTTP Server
// 	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
// 	handler := http.NewHandler(cache, l, cfg)

// 	v1.NewRouter(httpServer.Engine, handler, authUsecase, businessUsecase)

// 	httpServer.Start()

// 	// Waiting signal
// 	interrupt := make(chan os.Signal, 1)
// 	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

// 	select {
// 	case s := <-interrupt:
// 		l.Info("app - Run - signal: " + s.String())
// 	case err = <-httpServer.Notify():
// 		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
// 	}

// 	// Shutdown
// 	err = httpServer.Shutdown()
// 	if err != nil {
// 		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
// 	}
// }
