package service

// import (
// 	"context"

// 	v1 "sugurta/genproto/product_service"
// 	"sugurta/internal/entity"
// 	"sugurta/internal/usecase"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // OrderService is implementation of OrderServiceServer
// type OrderService struct {
// 	v1.UnimplementedOrderServiceServer
// 	useCase usecase.UseCase
// }

// // NewOrderService creates a new OrderService
// func NewOrderService(useCase usecase.UseCase) *OrderService {
// 	return &OrderService{
// 		useCase: useCase,
// 	}
// }

// // CreateOrder creates a new order
// func (s *OrderService) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.CreateOrderResponse, error) {
// 	if req.Order == nil {
// 		return nil, status.Error(codes.InvalidArgument, "order is required")
// 	}

// 	orderID, err := s.useCase.CreateOrder(ctx, entity.CreateOrderRequest{
// 		ClientID:   req.Order.UserId,
// 		BusinessID: req.Order.BusinessId,
// 		Status:     req.Order.Status,
// 		TotalCost:  int(req.Order.TotalAmount),
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.CreateOrderResponse{
// 		Id: *orderID,
// 	}, nil
// }

// // GetOrder gets an order by ID
// func (s *OrderService) GetOrder(ctx context.Context, req *v1.GetOrderRequest) (*v1.Order, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	order, err := s.useCase.GetOrder(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Order{
// 		Id:          order.ID,
// 		UserId:      order.ClientID,
// 		BusinessId:  order.BusinessID,
// 		TotalAmount: float64(order.TotalCost),
// 		Status:      order.Status,
// 		CreatedAt    time.Time:   order.CreatedAt    time.Time,
// 		UpdatedAt    time.Time:   order.UpdatedAt    time.Time,
// 	}, nil
// }

// // UpdateOrder updates an order
// func (s *OrderService) UpdateOrder(ctx context.Context, req *v1.UpdateOrderRequest) (*v1.Empty, error) {
// 	if req.Order == nil {
// 		return nil, status.Error(codes.InvalidArgument, "order is required")
// 	}

// 	err := s.useCase.UpdateOrder(ctx, entity.Order{
// 		ID:         req.Order.Id,
// 		ClientID:   req.Order.UserId,
// 		BusinessID: req.Order.BusinessId,
// 		TotalCost:  int(req.Order.TotalAmount),
// 		Status:     req.Order.Status,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // DeleteOrder deletes an order
// func (s *OrderService) DeleteOrder(ctx context.Context, req *v1.DeleteOrderRequest) (*v1.Empty, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	err := s.useCase.DeleteOrder(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }
