package controllers

import (
	"context"
	"database/sql"

	"github.com/philip-bui/grpc-errors"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/philip-bui/space-service/services/auth"
	"github.com/philip-bui/space-service/services/postgres"
)

func (s *Server) PostNew(ctx context.Context, req *pb.PostNewRequest) (*pb.PostNewResponse, error) {
	_, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}
	// TODO:
	return nil, nil
}

func (s *Server) PostGet(ctx context.Context, q *pb.ID) (*pb.Post, error) {
	post, err := postgres.GetPost(q.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return post, nil
}

func (s *Server) PostUpvote(ctx context.Context, q *pb.ID) (*pb.Empty, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := postgres.PostUpvote(userID, q.Id); err != nil {
		return nil, errors.ErrInternalServer
	}
	return &pb.Empty{}, nil
}

func (s *Server) PostUnUpvote(ctx context.Context, q *pb.ID) (*pb.Empty, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := postgres.PostUnUpvote(userID, q.Id); err != nil {
		return nil, errors.ErrInternalServer
	}
	return Empty, nil
}

func (s *Server) PostDownvote(ctx context.Context, q *pb.ID) (*pb.Empty, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := postgres.PostDownvote(userID, q.Id); err != nil {
		return nil, errors.ErrInternalServer
	}
	return Empty, nil
}

func (s *Server) PostUnDownvote(ctx context.Context, q *pb.ID) (*pb.Empty, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := postgres.PostUnDownvote(userID, q.Id); err != nil {
		return nil, errors.ErrInternalServer
	}
	return Empty, nil
}
