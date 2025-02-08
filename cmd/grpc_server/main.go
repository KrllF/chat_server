package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	res "github.com/KrllF/chat_server/cmd/db"
	"github.com/KrllF/chat_server/internal/config"
	desc "github.com/KrllF/chat_server/pkg/chat_server_v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatServerV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponce, error) {
	createresp, err := res.CreateChat(ctx, s.pool, req.GetUsernames())
	if err != nil {
		return nil, errors.New("error in creating chat")
	}
	return createresp, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendRequest) (*emptypb.Empty, error) {
	sendresp, err := res.SendMessage(ctx, s.pool, req.GetFrom(), req.GetText())
	if err != nil {
		return nil, errors.New("error in send message")
	}
	return sendresp, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatServerV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
