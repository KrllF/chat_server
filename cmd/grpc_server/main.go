package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/KrllF/chat_server/pkg/chat_server_v1"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatServerV1Server
}

func (s *server) Create(context.Context, *desc.CreateRequest) (*desc.CreateResponce, error) {
	return &desc.CreateResponce{}, nil
}

func (s *server) Delete(context.Context, *desc.DeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) SendMessage(context.Context, *desc.SendRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatServerV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
