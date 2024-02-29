package main

import (
	"context"
	"flag"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/kirillmc/chat-server/internal"
	"github.com/kirillmc/chat-server/internal/env"
	desc "github.com/kirillmc/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatV1Server
	p *pgxpool.Pool
}

// Create|Delete|SendMessage|
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	buildInsertChat := sq.Insert("chats").PlaceholderFormat(sq.Dollar).Columns("id").Values(sq.Expr("DEFAULT")).Suffix("RETURNING id")
	query, args, err := buildInsertChat.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}
	var chatID int64
	err = s.p.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Fatalf("failed to insert chat: %v", err)
	}

	buildInsertUsers := sq.Insert("chats_users").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_name")
	for _, elem := range req.Usernames {
		buildInsertUsers = buildInsertUsers.Values(chatID, elem)
	}
	query, args, err = buildInsertUsers.ToSql()
	if err != nil {
		log.Fatalf("failed to build query for chats_users: %v", err)
	}
	_, err = s.p.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to insert chats_users in databasse: %v", err)
	}
	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	builderDelete := sq.Delete("chats").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": req.GetId()})
	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatalf("failed to build DELETE query: %v", err)
	}
	_, err = s.p.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("fieled to delete chat: %v", err)
	}
	return nil, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	builderInsertMessage := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "from_user", "text", "timestamp").
		Values(req.ChatId, req.From, req.Text, req.Timestamp.AsTime())
	query, args, err := builderInsertMessage.ToSql()
	if err != nil {
		log.Fatalf("failed to build INSERT query to messages table: %v", err)
	}
	_, err = s.p.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to insert new message: %v", err)
	}
	return nil, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	//Считываем перемменные Вэнсдэй
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("fauled to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	// Создание пула соедениней с БД
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database")
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterChatV1Server(s, &server{p: pool})
	log.Printf("server is listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
