package converter

import (
	"github.com/kirillmc/chat-server/internal/model"
	desc "github.com/kirillmc/chat-server/pkg/chat_v1"
)

func ToChatFromCreateRequest(chat *desc.CreateRequest) *model.Chat {
	return &model.Chat{
		Usernames: chat.Usernames,
	}
}

func ToMessageFromSendMessageRequest(message *desc.SendMessageRequest) *model.Message {
	return &model.Message{
		ChatId:   message.Message.ChatId,
		UserFrom: message.Message.From,
		Text:     message.Message.Text,
	}
}

func ToUserChatConnectionFromConnectRequest(connectReq *desc.ConnectRequest) *model.UserChatConnection {
	return &model.UserChatConnection{
		ChatId:   connectReq.Id,
		UserName: connectReq.Username,
	}
}
