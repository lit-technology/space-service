package controllers

import (
	"context"

	"github.com/philip-bui/grpc-errors"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/philip-bui/space-service/services/auth"
	"github.com/philip-bui/space-service/services/postgres"
)

func (s *Server) PostComments(ctx context.Context, q *pb.IDQuery) (*pb.Comments, error) {
	comments, err := postgres.GetCommentsForPost(q.Id, q.Offset)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	// TODO:
	return &pb.Comments{Comments: comments}, nil
}

func (s *Server) CommentNew(ctx context.Context, req *pb.CommentNewRequest) (*pb.CommentNewResponse, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}
	commentID, err := postgres.InsertComment(req.PostID, userID, req.ReplyID, req.Text)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return &pb.CommentNewResponse{Id: commentID}, nil
}

func (s *Server) CommentUpvote(ctx context.Context, req *pb.ID) (*pb.Empty, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}
	if err := postgres.CommentUpvote(userID, req.Id); err != nil {
		return nil, errors.ErrInternalServer
	}
	return Empty, nil
}

func (s *Server) CommentUnUpvote(ctx context.Context, req *pb.ID) (*pb.Empty, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}
	if err := postgres.CommentUnUpvote(userID, req.Id); err != nil {
		return nil, errors.ErrInternalServer
	}
	return Empty, nil
}

func (s *Server) CommentDownvote(ctx context.Context, req *pb.ID) (*pb.Empty, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := postgres.CommentDownvote(userID, req.Id); err != nil {
		return nil, errors.ErrInternalServer
	}
	return Empty, nil
}

func (s *Server) CommentUnDownvote(ctx context.Context, req *pb.ID) (*pb.Empty, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := postgres.CommentUnDownvote(userID, req.Id); err != nil {
		return nil, errors.ErrInternalServer
	}
	return Empty, nil
}
