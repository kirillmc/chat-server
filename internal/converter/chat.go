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
		ChatId:   message.ChatId,
		UserFrom: message.From,
		Text:     message.Text,
	}
}
