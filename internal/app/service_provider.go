package app

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	descAccess "github.com/kirillmc/auth/pkg/access_v1"
	"github.com/kirillmc/chat-server/internal/api/chat"
	"github.com/kirillmc/chat-server/internal/client/rpc"
	"github.com/kirillmc/chat-server/internal/client/rpc/access"
	"github.com/kirillmc/chat-server/internal/config"
	"github.com/kirillmc/chat-server/internal/config/env"
	"github.com/kirillmc/chat-server/internal/interceptor"
	"github.com/kirillmc/chat-server/internal/repository"
	chatRepo "github.com/kirillmc/chat-server/internal/repository/chat"
	"github.com/kirillmc/chat-server/internal/service"
	chatService "github.com/kirillmc/chat-server/internal/service/chat"
	"github.com/kirillmc/platform_common/pkg/closer"
	"github.com/kirillmc/platform_common/pkg/db"
	"github.com/kirillmc/platform_common/pkg/db/pg"
)

// содержит все зависимости, необходимые в рамках приложения
type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig
	accesConfig   config.AccessConfig

	accessClient      rpc.AccessClient
	dbClient          db.Client
	interceptorClient *interceptor.Interceptor

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

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		httpConfig, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %v", err)
		}

		s.httpConfig = httpConfig
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		swaggerConfig, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %v", err)
		}

		s.swaggerConfig = swaggerConfig
	}

	return s.swaggerConfig
}

func (s *serviceProvider) AccessConfig() config.AccessConfig {
	if s.accesConfig == nil {
		cfg, err := env.NewAccessConfig()
		if err != nil {
			log.Fatalf("failed to get acces configs: %v", err)
		}

		s.accesConfig = cfg
	}

	return s.accesConfig
}

func (s *serviceProvider) AccessClient() rpc.AccessClient {
	if s.accessClient == nil {
		cfg := s.AccessConfig()
		creds, err := credentials.NewClientTLSFromFile(cfg.CertPath(), "")
		if err != nil {
			log.Fatalf("failed to get credentials of access: %v", err)
		}

		conn, err := grpc.Dial(
			cfg.Address(),
			grpc.WithTransportCredentials(creds),
		)
		if err != nil {
			log.Fatalf("failed to connect to access: %v", err)
		}

		s.accessClient = access.NewAccessClient(descAccess.NewAccessV1Client(conn))
	}

	return s.accessClient
}

func (s *serviceProvider) InterceptorClient() *interceptor.Interceptor {
	if s.interceptorClient == nil {
		s.interceptorClient = &interceptor.Interceptor{
			client: s.AccessClient(),
		}
	}

	return s.interceptorClient
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
