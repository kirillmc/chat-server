package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/kirillmc/chat-server/internal/model"
	"github.com/kirillmc/chat-server/internal/repository"
	repositoryMocks "github.com/kirillmc/chat-server/internal/repository/mocks"
	"github.com/kirillmc/chat-server/internal/service/chat"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		req *model.Message
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatId   = gofakeit.Int64()
		userFrom = gofakeit.Name()
		text     = gofakeit.BeerAlcohol()

		repositoryErr = fmt.Errorf("error of repository layer")

		modelMessage = &model.Message{
			ChatId:   chatId,
			UserFrom: userFrom,
			Text:     text,
		}
	)

	tests := []struct {
		name               string
		args               args
		err                error
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "success send_message_from case",
			args: args{
				ctx: ctx,
				req: modelMessage,
			},
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, modelMessage).Return(nil)
				return mock
			},
		},
		{
			name: "service error send_message_from",
			args: args{
				ctx: ctx,
				req: modelMessage,
			},
			err: repositoryErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, modelMessage).Return(repositoryErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepositoryMock := tt.chatRepositoryMock(mc)
			service := chat.NewService(chatRepositoryMock)

			err := service.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
