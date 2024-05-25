package access

import (
	"context"

	descAccess "github.com/kirillmc/auth/pkg/access_v1"
	"github.com/kirillmc/chat-server/internal/client/rpc"
)

type accessClient struct {
	client descAccess.AccessV1Client
}

var _ rpc.AccessClient = (*accessClient)(nil)

func NewAccessClient(client descAccess.AccessV1Client) rpc.AccessClient {
	return &accessClient{
		client: client,
	}
}

func (c *accessClient) Check(ctx context.Context, endpoint string) error {
	_, err := c.client.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: endpoint,
	})
	return err
}
