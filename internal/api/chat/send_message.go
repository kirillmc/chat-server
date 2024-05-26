package chat

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kirillmc/chat-server/internal/converter"
	desc "github.com/kirillmc/chat-server/pkg/chat_v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetMessage().GetChatId()]
	i.mxChannel.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}
	chatChan <- &desc.ConnectResponse{
		Message: req.GetMessage(),
	}

	err := i.chatService.SendMessage(ctx, converter.ToMessageFromSendMessageRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
