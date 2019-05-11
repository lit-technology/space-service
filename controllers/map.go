package controllers

import (
	"context"

	pb "github.com/philip-bui/space-service/protos"
	"github.com/philip-bui/space-service/services/auth"
	//"github.com/philip-bui/space-service/services/postgres"
)

func (s *Server) Map(ctx context.Context, req *pb.MapRequest) (*pb.MapResponse, error) {
	claims, err := auth.GetClaimsFromContext(ctx)
	if err != nil {
		return GetMapSpaces(ctx, req)
	} else {
		return GetMapSpacesForUser(ctx, req, claims)
	}
}

func GetMapSpaces(ctx context.Context, req *pb.MapRequest) (*pb.MapResponse, error) {
	// TODO
	return nil, nil
}

func GetMapSpacesForUser(ctx context.Context, req *pb.MapRequest, claims map[string]interface{}) (*pb.MapResponse, error) {
	// TODO
	return nil, nil
}
