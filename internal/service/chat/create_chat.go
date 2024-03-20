package chat

import (
	"context"

	"github.com/kirillmc/chat-server/internal/model"
)

func (s *serv) CreateChat(ctx context.Context, req *model.Chat) (int64, error) {
	id, err := s.chatRepository.CreateChat(ctx, req)
	if err != nil {
		return 0, err
	}

	err = s.chatRepository.AddUsers(ctx, req.Usernames, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
