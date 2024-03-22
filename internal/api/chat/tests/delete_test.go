package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/kirillmc/chat-server/internal/api/chat"
	"github.com/kirillmc/chat-server/internal/service"
	serviceMocks "github.com/kirillmc/chat-server/internal/service/mocks"
	desc "github.com/kirillmc/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("error of service layer")

		req = &desc.DeleteRequest{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success deletee case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.DeleteChatMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "service error deletee",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.DeleteChatMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			api := chat.NewImplementation(chatServiceMock)

			_, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
