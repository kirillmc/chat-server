package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/kirillmc/chat-server/internal/client/rpc"
)

type Client struct {
	Client rpc.AccessClient
}

func (c *Client) PolicyInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadate is not provided")
	}

	err := c.Client.Check(metadata.NewOutgoingContext(ctx, md), info.FullMethod)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
