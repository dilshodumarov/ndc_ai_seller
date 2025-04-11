package grpc

// import (
// 	"context"
// 	"fmt"
// 	"net"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"sugurta/internal/pkg/config"
// 	auth "sugurta/genproto/auth_service"
// 	business "sugurta/genproto/business_service"
// 	product "sugurta/genproto/product_service"
// 	"sugurta/internal/grpc/service"
// 	"sugurta/internal/repo"
// 	"sugurta/internal/repo/auth"
// 	"sugurta/internal/usecase"
// 	"sugurta/internal/pkg/logger"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"
// )

// // Server is gRPC server
// type Server struct {
// 	cfg        *config.Config
// 	logger     logger.Interface
// 	authRepo auth.
// 	useCase    usecase.UseCase
// 	grpcServer *grpc.Server
// }

// // NewServer creates a new gRPC server
// func NewServer(cfg *config.Config, logger logger.Interface, storage *repo.Storage, useCase usecase.UseCase) *Server {
// 	return &Server{
// 		cfg:     cfg,
// 		logger:  logger,
// 		storage: &repo.Storage{

// 		},
// 		useCase: useCase,
// 	}
// }

// // Run runs gRPC server
// func (s *Server) Run() error {
// 	// Create a listener on TCP port
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.GRPC.Port))
// 	if err != nil {
// 		s.logger.Error("failed to listen", logger.Any("error", err))
// 		return err
// 	}

// 	// Create a gRPC server
// 	s.grpcServer = grpc.NewServer()

// 	// Register reflection service on gRPC server
// 	reflection.Register(s.grpcServer)

// 	// Register all services
// 	s.registerServices()

// 	// Start gRPC server
// 	go func() {
// 		s.logger.Info("Starting gRPC server", logger.String("port", s.cfg.GRPC.Port))
// 		if err := s.grpcServer.Serve(lis); err != nil {
// 			s.logger.Fatal("failed to serve", logger.Any("error", err))
// 		}
// 	}()

// 	// Wait for interrupt signal to gracefully shutdown the server
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	// Gracefully stop server
// 	s.grpcServer.GracefulStop()
// 	s.logger.Info("gRPC server shutting down...")

// 	return nil
// }

// // registerServices registers all gRPC services
// func (s *Server) registerServices() {
// 	// Register User service
// 	userService := service.NewUserService(s.useCase)
// 	auth.RegisterUserServiceServer(s.grpcServer, userService)

// 	// Register Role service
// 	roleService := service.NewRoleService(s.useCase)
// 	auth.RegisterRoleServiceServer(s.grpcServer, roleService)

// 	// Register ClientType service
// 	clientTypeService := service.NewClientTypeService(s.useCase)
// 	auth.RegisterClientTypeServiceServer(s.grpcServer, clientTypeService)

// 	// Register Business service
// 	businessService := service.NewBusinessService(s.useCase)
// 	business.RegisterBusinessServiceServer(s.grpcServer, businessService)

// 	// Register Product service
// 	productService := service.NewProductService(s.useCase)
// 	product.RegisterProductServiceServer(s.grpcServer, productService)

// 	// Register Category service
// 	categoryService := service.NewCategoryService(s.useCase)
// 	product.RegisterCategoryServiceServer(s.grpcServer, categoryService)

// 	// Register Order service
// 	orderService := service.NewOrderService(s.useCase)
// 	product.RegisterOrderServiceServer(s.grpcServer, orderService)

// 	// Register OrderProducts service
// 	orderProductsService := service.NewOrderProductsService(s.useCase)
// 	product.RegisterOrderProductsServiceServer(s.grpcServer, orderProductsService)
// }

// // Stop stops the gRPC server
// func (s *Server) Stop(ctx context.Context) error {
// 	s.grpcServer.GracefulStop()
// 	return nil
// }
