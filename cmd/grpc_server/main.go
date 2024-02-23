package main

import (
	"context"
	"fmt"
	desc "github.com/kirillmc/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const grpcPort = 127001

type server struct {
	desc.UnimplementedChatV1Server
}

// Create|Delete|SendMessage|
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	return &desc.CreateResponse{
		Id: 101,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) error {
	log.Printf("user_%d is deleted", req.GetId())
	return nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) error {
	log.Printf("message was sent at %v", req.GetTimestamp())
	return nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterChatV1Server(s, &server{})
	log.Printf("server is listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
