package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kirillmc/chat-server/internal/converter"
	desc "github.com/kirillmc/chat-server/pkg/chat_v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.chatService.SendMessage(ctx, converter.ToMessageFromSendMessageRequest(req))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
