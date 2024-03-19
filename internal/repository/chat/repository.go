package chat

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/kirillmc/chat-server/internal/client/db"
	"github.com/kirillmc/chat-server/internal/model"
	"github.com/kirillmc/chat-server/internal/repository"
)

const (
	chatsTable      = "chats"
	chatsUsersTable = "chats_users"
	messagesTable   = "messages"
	idColumn        = "id"
	chatIdColumn    = "chat_id"
	userNameColumn  = "user_name"
	fromUserColumn  = "from_user"
	textColumn      = "text"

	defaultText     = "DEFAULT"
	returningIDText = "RETURNING id"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{
		db: db,
	}
}

// Create|Delete|SendMessage|
func (r *repo) CreateChat(ctx context.Context, req *model.Chat) (int64, error) {
	buildInsertChat := sq.Insert(chatsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(idColumn).
		Values(sq.Expr(defaultText)).
		Suffix(returningIDText)

	query, args, err := buildInsertChat.ToSql()
	if err != nil {
		return 0, err
	}

	var chatID int64

	q := db.Query{
		Name:     "chat_repository.CreateChat",
		QueryRaw: query,
	}
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (r *repo) AddUser(ctx context.Context, userName string, chatId int64) error {
	buildInsertUsers := sq.Insert(chatsUsersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIdColumn, userNameColumn)

	buildInsertUsers = buildInsertUsers.Values(chatId, userName)

	query, args, err := buildInsertUsers.ToSql()

	q := db.Query{
		Name:     "chat_repository.AddUsersToNewChat",
		QueryRaw: query,
	}
	if err != nil {
		return err
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteChat(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(chatsTable).PlaceholderFormat(sq.Dollar).Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.DeleteChat",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) SendMessage(ctx context.Context, req *model.Message) error {
	builderInsertMessage := sq.Insert(messagesTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIdColumn, fromUserColumn, textColumn).
		Values(req.ChatId, req.UserFrom, req.Text)

	query, args, err := builderInsertMessage.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.SendMessage",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
