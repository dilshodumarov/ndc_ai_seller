package grpc_service_clients

// import (
// 	"fmt"

// 	"evrone_api_gateway/genproto/content_service"
// 	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
// 	"google.golang.org/grpc"

// 	"evrone_api_gateway/internal/pkg/config"
// )

// type ServiceClient interface {
// 	ContentService() content_service.ContentServiceClient
// 	Close()
// }

// type serviceClient struct {
// 	connections    []*grpc.ClientConn
// 	contentService content_service.ContentServiceClient
// }

// func New(cfg *config.Config) (ServiceClient, error) {
// 	connContentService, err := grpc.Dial(
// 		fmt.Sprintf("%s%s", cfg.ContentService.Host, cfg.ContentService.Port),
// 		grpc.WithInsecure(),
// 		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
// 		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &serviceClient{
// 		contentService: content_service.NewContentServiceClient(connContentService),
// 		connections: []*grpc.ClientConn{
// 			connContentService,
// 		},
// 	}, nil
// }

// func (s *serviceClient) ContentService() content_service.ContentServiceClient {
// 	return s.contentService
// }

// func (s *serviceClient) Close() {
// 	for _, conn := range s.connections {
// 		if err := conn.Close(); err != nil {
// 			// should be replaced by logger soon
// 			fmt.Printf("error while closing grpc connection: %v", err)
// 		}
// 	}
// }
