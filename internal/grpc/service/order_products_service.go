package service

// import (
// 	"context"

// 	v1 "sugurta/genproto/product_service"
// 	"sugurta/internal/entity"
// 	"sugurta/internal/usecase"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // OrderProductsService is implementation of OrderProductsServiceServer
// type OrderProductsService struct {
// 	v1.UnimplementedOrderProductsServiceServer
// 	useCase usecase.UseCase
// }

// // NewOrderProductsService creates a new OrderProductsService
// func NewOrderProductsService(useCase usecase.UseCase) *OrderProductsService {
// 	return &OrderProductsService{
// 		useCase: useCase,
// 	}
// }

// // CreateOrderProducts creates a new order products entry
// func (s *OrderProductsService) CreateOrderProducts(ctx context.Context, req *v1.CreateOrderProductsRequest) (*v1.Empty, error) {
// 	if req.OrderProducts == nil {
// 		return nil, status.Error(codes.InvalidArgument, "order_products is required")
// 	}

// 	err := s.useCase.CreateOrderProducts(ctx, entity.OrderProducts{
// 		ID:        req.OrderProducts.Id,
// 		OrderID:   req.OrderProducts.OrderId,
// 		ProductID: req.OrderProducts.ProductId,
// 		Count:     int(req.OrderProducts.Quantity),
// 		Cost:      int(req.OrderProducts.Quantity),
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // GetOrderProducts gets order products by ID
// func (s *OrderProductsService) GetOrderProducts(ctx context.Context, req *v1.GetOrderProductsRequest) (*v1.OrderProducts, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	orderProducts, err := s.useCase.GetOrderProducts(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	result := &v1.OrderProducts{
// 		Id:        orderProducts.ID,
// 		OrderId:   orderProducts.OrderID,
// 		ProductId: orderProducts.ProductID,
// 		Quantity:  int32(orderProducts.Count),
// 		Price:     float64(orderProducts.Cost),
// 		CreatedAt    time.Time: orderProducts.CreatedAt    time.Time,
// 		UpdatedAt    time.Time: orderProducts.UpdatedAt    time.Time,
// 	}

// 	if orderProducts.ProductID != "" {
// 		result.Product = &v1.Product{
// 			Id:    orderProducts.ProductID,
// 			Price: float64(orderProducts.Cost),
// 		}

// 		// if orderProducts.Category.ID != "" {
// 		// 	result.Product.Category = &v1.Category{
// 		// 		Id:          orderProducts.Product.Category.ID,
// 		// 		Name:        orderProducts.Product.Category.Name,
// 		// 		Description: orderProducts.Product.Category.Description,
// 		// 	}
// 		// }
// 	}

// 	return result, nil
// }

// // UpdateOrderProducts updates order products
// func (s *OrderProductsService) UpdateOrderProducts(ctx context.Context, req *v1.UpdateOrderProductsRequest) (*v1.Empty, error) {
// 	if req.OrderProducts == nil {
// 		return nil, status.Error(codes.InvalidArgument, "order_products is required")
// 	}

// 	err := s.useCase.UpdateOrderProducts(ctx, entity.OrderProducts{
// 		ID:        req.OrderProducts.Id,
// 		OrderID:   req.OrderProducts.OrderId,
// 		ProductID: req.OrderProducts.ProductId,
// 		Count:     int(req.OrderProducts.Quantity),
// 		Cost:      int(req.OrderProducts.Price),
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // DeleteOrderProducts deletes order products
// func (s *OrderProductsService) DeleteOrderProducts(ctx context.Context, req *v1.DeleteOrderProductsRequest) (*v1.Empty, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	err := s.useCase.DeleteOrderProducts(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }
