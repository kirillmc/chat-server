package app

import (
	"context"
	"log"

	"github.com/kirillmc/chat-server/internal/api/chat"
	"github.com/kirillmc/chat-server/internal/client/db"
	"github.com/kirillmc/chat-server/internal/client/db/pg"
	"github.com/kirillmc/chat-server/internal/closer"
	"github.com/kirillmc/chat-server/internal/config"
	"github.com/kirillmc/chat-server/internal/config/env"
	"github.com/kirillmc/chat-server/internal/repository"
	chatRepo "github.com/kirillmc/chat-server/internal/repository/chat"
	"github.com/kirillmc/chat-server/internal/service"
	chatService "github.com/kirillmc/chat-server/internal/service/chat"
)

// содержит все зависимости, необходимые в рамках приложения
type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient db.Client

	chatRepository repository.ChatRepository
	chatService    service.ChatService

	chatImpl *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// если Get - в GO Get НЕ УКАЗЫВАЮТ: НЕ GetPGConfig, A PGConfig
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		pgConfig, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err) // делаем Log.Fatalf, чтобы не обрабатывать ошибку в другом месте
			// + инициализация происходит при старте приложения, поэтому если ошибка - можно и сервер уронить
			// можно кинуть panic()
		}
		s.pgConfig = pgConfig
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		grpcConfig, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		s.grpcConfig = grpcConfig
	}
	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("oing error: %s", err.Error())
		}

		closer.Add(cl.Close)
		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepo.NewRepository(s.DBClient(ctx))
	}
	return s.chatRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(s.UserRepository(ctx))
	}
	return s.chatService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *chat.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewImplementation(s.UserService(ctx))
	}
	return s.chatImpl
}
