package client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/go-mockingcode/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthGRPCClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthGRPCClient(grpcURL string) (*AuthGRPCClient, error) {
	slog.Info("connecting to auth gRPC service", slog.String("url", grpcURL))

	conn, err := grpc.NewClient(grpcURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthServiceClient(conn)

	return &AuthGRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *AuthGRPCClient) ValidateToken(token string) (*ValidateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		slog.Error("grpc call failed", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Debug("grpc auth response",
		slog.Bool("valid", resp.Valid),
		slog.Int64("user_id", resp.UserId),
	)

	return &ValidateResponse{
		Valid:  resp.Valid,
		UserID: fmt.Sprintf("%d", resp.UserId), // Convert int64 to string
	}, nil
}

func (c *AuthGRPCClient) Close() error {
	return c.conn.Close()
}
