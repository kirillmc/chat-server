package chat

import (
	"sync"

	"github.com/kirillmc/chat-server/internal/service"
	desc "github.com/kirillmc/chat-server/pkg/chat_v1"
)

type Chat struct {
	streams map[string]desc.ChatV1_ConnectServer
	m       sync.RWMutex
}

type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService

	chats  map[int64]*Chat
	mxChat sync.RWMutex

	channels  map[int64]chan *desc.ConnectResponse
	mxChannel sync.RWMutex
}

func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
		chats:       make(map[int64]*Chat),
		channels:    make(map[int64]chan *desc.ConnectResponse),
	}
}
