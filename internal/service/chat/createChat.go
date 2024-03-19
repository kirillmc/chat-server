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

	for _, elem := range req.Usernames {
		err = s.chatRepository.AddUser(ctx, elem, id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}
