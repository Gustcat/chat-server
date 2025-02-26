package interceptor

import (
	"context"
	"fmt"
	"github.com/Gustcat/auth/pkg/access_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

func AuthInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, err := grpc.NewClient(
		"auth-backend-1:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("error creating gRPC client: %v", err)
		return nil, err
	}

	client := access_v1.NewAccessV1Client(conn)

	_, err = client.Check(ctx, &access_v1.CheckRequest{
		EndpointAddress: info.FullMethod,
	})
	if err != nil {
		log.Printf("error sending grpc request: %v", err)
		return nil, err
	}

	return handler(ctx, req)
}
