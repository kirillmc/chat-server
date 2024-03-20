package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/kirillmc/chat-server/pkg/chat_v1"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.DeleteChat(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}
