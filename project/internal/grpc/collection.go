package grpc

import (
	"context"
	"encoding/json"
	"log/slog"

	pb "github.com/go-mockingcode/proto"
)

// GetCollectionSchema returns collection schema by name (for optional validation)
func (s *ProjectGRPCServer) GetCollectionSchema(ctx context.Context, req *pb.GetCollectionSchemaRequest) (*pb.GetCollectionSchemaResponse, error) {
	slog.Debug("grpc: getting collection schema",
		slog.Int64("project_id", req.ProjectId),
		slog.String("collection_name", req.CollectionName),
	)

	// Получаем коллекцию по имени
	collection, err := s.collectionService.GetCollectionByName(req.ProjectId, req.CollectionName)
	if err != nil {
		slog.Error("grpc: failed to get collection", slog.String("error", err.Error()))
		return &pb.GetCollectionSchemaResponse{
			Found: false,
		}, nil
	}

	// Коллекция не найдена - это нормально для гибридного подхода
	if collection == nil {
		slog.Debug("grpc: collection schema not found (schema-less mode)")
		return &pb.GetCollectionSchemaResponse{
			Found: false,
		}, nil
	}

	// Конвертируем fields в JSON string
	fieldsJSON, err := json.Marshal(collection.Fields)
	if err != nil {
		slog.Error("grpc: failed to marshal fields", slog.String("error", err.Error()))
		return &pb.GetCollectionSchemaResponse{
			Found: false,
		}, nil
	}

	slog.Debug("grpc: collection schema found",
		slog.Int64("collection_id", collection.ID),
		slog.String("name", collection.Name),
		slog.Int("fields_count", len(collection.Fields)),
	)

	return &pb.GetCollectionSchemaResponse{
		Found:        true,
		CollectionId: collection.ID,
		Name:         collection.Name,
		Description:  collection.Description,
		FieldsJson:   string(fieldsJSON),
		IsActive:     collection.IsActive,
	}, nil
}

