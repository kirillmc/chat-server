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

func TestCreate(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		req *model.Chat
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		users = []string{gofakeit.BeerName(), gofakeit.BeerName(), gofakeit.BeerName()}

		repositoryErr = fmt.Errorf("error of repository layer")

		req = &model.Chat{
			Usernames: users,
		}
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "success create case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx, req).Return(id, nil)
				mock.AddUsersMock.Expect(ctx, req.Usernames, id).Return(nil)
				return mock
			},
		},
		{
			name: "repository error of add users",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repositoryErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx, req).Return(id, nil)
				mock.AddUsersMock.Expect(ctx, req.Usernames, id).Return(repositoryErr)
				return mock
			},
		},

		{
			name: "repository error create",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repositoryErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx, req).Return(0, repositoryErr)
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

			newId, err := service.CreateChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newId)
		})
	}
}
