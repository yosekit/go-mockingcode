package client

import (
	"context"
	"log/slog"
	"time"

	pb "github.com/go-mockingcode/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProjectGRPCClient struct {
	client pb.ProjectServiceClient
	conn   *grpc.ClientConn
}

func NewProjectGRPCClient(grpcURL string) (*ProjectGRPCClient, error) {
	slog.Info("connecting to project gRPC service", slog.String("url", grpcURL))

	conn, err := grpc.NewClient(grpcURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewProjectServiceClient(conn)

	return &ProjectGRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

type ProjectInfo struct {
	Valid       bool
	ProjectID   int64
	ProjectName string
	UserID      int64
}

func (c *ProjectGRPCClient) ValidateAPIKey(apiKey string) (*ProjectInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.ValidateAPIKey(ctx, &pb.ValidateAPIKeyRequest{
		ApiKey: apiKey,
	})
	if err != nil {
		slog.Error("grpc call failed", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Debug("grpc project response",
		slog.Bool("valid", resp.Valid),
		slog.Int64("project_id", resp.ProjectId),
	)

	return &ProjectInfo{
		Valid:       resp.Valid,
		ProjectID:   resp.ProjectId,
		ProjectName: resp.ProjectName,
		UserID:      resp.UserId,
	}, nil
}

type CollectionSchema struct {
	Found       bool
	CollectionID int64
	Name        string
	Description string
	FieldsJSON  string
	IsActive    bool
}

func (c *ProjectGRPCClient) GetCollectionSchema(projectID int64, collectionName string) (*CollectionSchema, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetCollectionSchema(ctx, &pb.GetCollectionSchemaRequest{
		ProjectId:      projectID,
		CollectionName: collectionName,
	})
	if err != nil {
		slog.Error("grpc GetCollectionSchema failed", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Debug("grpc collection schema response",
		slog.Bool("found", resp.Found),
		slog.String("collection", collectionName),
	)

	return &CollectionSchema{
		Found:        resp.Found,
		CollectionID: resp.CollectionId,
		Name:         resp.Name,
		Description:  resp.Description,
		FieldsJSON:   resp.FieldsJson,
		IsActive:     resp.IsActive,
	}, nil
}

func (c *ProjectGRPCClient) Close() error {
	return c.conn.Close()
}
