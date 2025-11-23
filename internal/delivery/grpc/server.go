package grpc

import (
	"context"

	"github.com/Snake1-1eyes/auth-service-it/internal/usecase/auth"
	pb "github.com/Snake1-1eyes/auth-service-it/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	uc *auth.UseCase
}

func NewServer(uc *auth.UseCase) *Server {
	return &Server{uc: uc}
}

func (s *Server) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	userID, err := s.uc.SignUp(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to sign up: %v", err)
	}

	return &pb.SignUpResponse{UserId: userID}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	sessionID, err := s.uc.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to login: %v", err)
	}

	return &pb.LoginResponse{SessionId: sessionID}, nil
}
