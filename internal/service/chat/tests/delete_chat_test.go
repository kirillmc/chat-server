package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/kirillmc/chat-server/internal/repository"
	repositoryMocks "github.com/kirillmc/chat-server/internal/repository/mocks"
	"github.com/kirillmc/chat-server/internal/service/chat"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repositoryErr = fmt.Errorf("error of repository layer")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "success deletee case",
			args: args{
				ctx: ctx,
				req: id,
			},
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.DeleteChatMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "repository error deletee",
			args: args{
				ctx: ctx,
				req: id,
			},
			err: repositoryErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.DeleteChatMock.Expect(ctx, id).Return(repositoryErr)
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

			err := service.DeleteChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
