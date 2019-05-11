package controllers

import (
	"context"

	"github.com/philip-bui/grpc-errors"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/philip-bui/space-service/services/auth"
	//"github.com/philip-bui/space-service/services/postgres"
)

func (s *Server) Home(ctx context.Context, req *pb.HomeRequest) (*pb.HomeResponse, error) {
	claims, err := auth.GetClaimsFromContext(ctx)
	if err != nil {
		return GetHomePosts(ctx)
	} else {
		return GetHomePostsForUser(ctx, claims)
	}
}

func GetHomePosts(ctx context.Context) (*pb.HomeResponse, error) {
	// TODO
	return nil, errors.ErrInternalServer
}

func GetHomePostsForUser(ctx context.Context, claims map[string]interface{}) (*pb.HomeResponse, error) {
	// TODO
	return nil, nil
}
