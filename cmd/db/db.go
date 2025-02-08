package db

import (
	"context"
	"fmt"
	"log"
	"strings"

	desc "github.com/KrllF/chat_server/pkg/chat_server_v1"
	"google.golang.org/protobuf/types/known/emptypb"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateChat(ctx context.Context, pool *pgxpool.Pool, usernames []string) (*desc.CreateResponce, error) {
	memb := strings.Join(usernames, "&")
	builderInsert := sq.Insert("chats").
		PlaceholderFormat(sq.Dollar).
		Columns("members").
		Values(memb).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Failed to build query: %v", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	var chatID int
	err = pool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Printf("Failed to insert chat: query=%s, args=%v, error=%v", query, args, err)
		return nil, fmt.Errorf("failed to insert chat: %w", err)
	}
	log.Printf("Inserted chat with id: %d", int64(chatID))
	return &desc.CreateResponce{Id: int64(chatID)}, nil
}

func SendMessage(ctx context.Context, pool *pgxpool.Pool, from string, text string) (*emptypb.Empty, error) {
	builderInsert := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("sender", "letter").
		Values(from, text)
	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}
	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed create message: %v", err)
		return nil, err
	}
	log.Printf("create %d message", res.RowsAffected())
	return &emptypb.Empty{}, nil
}
