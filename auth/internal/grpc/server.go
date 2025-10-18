package grpc

import (
	"context"
	"log/slog"

	"github.com/go-mockingcode/auth/internal/service"
	pb "github.com/go-mockingcode/proto"
)

type AuthGRPCServer struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthGRPCServer(authService *service.AuthService) *AuthGRPCServer {
	return &AuthGRPCServer{
		authService: authService,
	}
}

// ValidateToken validates JWT access token via gRPC
func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	slog.Debug("grpc: validating token", slog.String("token_length", string(rune(len(req.Token)))))

	user, err := s.authService.ValidateAccessToken(req.Token)
	if err != nil {
		slog.Warn("grpc: invalid token", slog.String("error", err.Error()))
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	slog.Debug("grpc: token valid",
		slog.Int64("user_id", user.ID),
		slog.String("email", user.Email),
	)

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: user.ID,
		Email:  user.Email,
	}, nil
}

