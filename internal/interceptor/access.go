package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/kirillmc/chat-server/internal/client/rpc"
	"github.com/opentracing/opentracing-go"
)

type Interceptor struct {
	Client rpc.AccessClient
}

func (i *Interceptor) PolicyInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "check")
	defer span.Finish()

	err := i.Client.Check(metadata.NewOutgoingContext(ctx, md), info.FullMethod)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}
