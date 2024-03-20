package service

import (
	"context"

	"github.com/kirillmc/chat-server/internal/model"
)

type ChatService interface {
	CreateChat(ctx context.Context, req *model.Chat) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, req *model.Message) error
}
