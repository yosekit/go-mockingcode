package grpc

import (
	"context"
	"log/slog"

	"github.com/go-mockingcode/project/internal/service"
	pb "github.com/go-mockingcode/proto"
)

type ProjectGRPCServer struct {
	pb.UnimplementedProjectServiceServer
	projectService    *service.ProjectService
	collectionService *service.CollectionService
}

func NewProjectGRPCServer(projectService *service.ProjectService, collectionService *service.CollectionService) *ProjectGRPCServer {
	return &ProjectGRPCServer{
		projectService:    projectService,
		collectionService: collectionService,
	}
}

// ValidateAPIKey validates project API key and returns project info via gRPC
func (s *ProjectGRPCServer) ValidateAPIKey(ctx context.Context, req *pb.ValidateAPIKeyRequest) (*pb.ValidateAPIKeyResponse, error) {
	// Mask API key for logging
	maskedKey := req.ApiKey
	if len(req.ApiKey) > 8 {
		maskedKey = req.ApiKey[:8] + "..."
	}
	slog.Debug("grpc: validating API key", slog.String("api_key", maskedKey))

	project, err := s.projectService.GetProjectByAPIKey(req.ApiKey)
	if err != nil || project == nil {
		slog.Warn("grpc: invalid API key")
		return &pb.ValidateAPIKeyResponse{
			Valid: false,
		}, nil
	}

	slog.Debug("grpc: API key valid",
		slog.Int64("project_id", project.ID),
		slog.String("project_name", project.Name),
	)

	return &pb.ValidateAPIKeyResponse{
		Valid:       true,
		ProjectId:   project.ID,
		ProjectName: project.Name,
		UserId:      project.UserID,
	}, nil
}

