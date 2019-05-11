package controllers

import (
	"context"

	"github.com/philip-bui/grpc-errors"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/philip-bui/space-service/services/postgres"
)

func (s *Server) PostSearch(ctx context.Context, q *pb.Query) (*pb.Posts, error) {
	// TODO:
	return nil, nil
}

func (s *Server) UserSearch(ctx context.Context, q *pb.Query) (*pb.Users, error) {
	users, err := postgres.SearchUser(q.Query)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return &pb.Users{Users: users}, nil
}

func (s *Server) TagSearch(ctx context.Context, q *pb.Query) (*pb.Tags, error) {
	tags, err := postgres.SearchTag(q.Query)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return &pb.Tags{Tags: tags}, nil
}

func (s *Server) SpaceSearch(ctx context.Context, q *pb.Query) (*pb.Spaces, error) {
	return nil, nil
}
