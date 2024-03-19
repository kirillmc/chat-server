package chat

import (
	"context"

	"github.com/kirillmc/chat-server/internal/converter"
	desc "github.com/kirillmc/chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.chatService.CreateChat(ctx, converter.ToChatFromCreateRequest(req))
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
