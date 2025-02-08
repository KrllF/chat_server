package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	res "github.com/KrllF/chat_server/cmd/db"
	desc "github.com/KrllF/chat_server/pkg/chat_server_v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	grpcPort = 50052
	dbDSN    = "host=localhost port=54320 dbname=chat-server user=chat-server-user password=chat-server-password sslmode=disable"
)

type server struct {
	desc.UnimplementedChatServerV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponce, error) {
	createresp, err := res.CreateChat(ctx, req.GetUsernames())
	if err != nil {
		return nil, errors.New("error in creating chat")
	}
	return createresp, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendRequest) (*emptypb.Empty, error) {
	sendresp, err := res.SendMessage(ctx, req.GetFrom(), req.GetText())
	if err != nil {
		return nil, errors.New("error in send message")
	}
	return sendresp, nil
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
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	res.DB = pool
	defer pool.Close()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
