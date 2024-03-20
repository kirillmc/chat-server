package model

type Chat struct {
	Usernames []string
}

type Message struct {
	ChatId   int64
	UserFrom string
	Text     string
}
