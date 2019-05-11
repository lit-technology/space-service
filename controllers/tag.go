package controllers

import (
	"context"
	"database/sql"

	"github.com/philip-bui/grpc-errors"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/philip-bui/space-service/services/postgres"
)

func (s *Server) TagGet(ctx context.Context, q *pb.StringID) (*pb.Tag, error) {
	tag, err := postgres.GetTag(q.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return tag, nil
}

func (s *Server) TagPosts(ctx context.Context, q *pb.StringIDQuery) (*pb.Posts, error) {
	posts, err := postgres.GetPostsForTag(q.Id, q.Offset)
	if err != nil {
		return nil, err
	}
	return &pb.Posts{Posts: posts}, nil
}
