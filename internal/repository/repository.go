package repository

import (
	"context"

	"github.com/kirillmc/chat-server/internal/model"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, req *model.Chat) (int64, error)
	AddUsers(ctx context.Context, userNames []string, chatId int64) error
	DeleteChat(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, req *model.Message) error
}
