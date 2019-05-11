package controllers

import (
	"context"

	"github.com/philip-bui/grpc-errors"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/philip-bui/space-service/services/auth"
	//	"github.com/philip-bui/space-service/services/postgres"
)

func (s *Server) UserProfileGet(ctx context.Context, _ *pb.Empty) (*pb.UserProfile, error) {
	_, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// TODO:
	return nil, errors.ErrInternalServer
}

func (s *Server) UserProfileEdit(ctx context.Context, req *pb.UserProfileRequest) (*pb.UserProfileResponse, error) {
	// TODO:
	return nil, errors.ErrInternalServer
}
