package client

// TokenValidator is an interface for token validation
// Implemented by both HTTP and gRPC auth clients
type TokenValidator interface {
	ValidateToken(token string) (*ValidateResponse, error)
}

