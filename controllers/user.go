package controllers

import (
	"context"
	"database/sql"

	"github.com/philip-bui/grpc-errors"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/philip-bui/space-service/services/postgres"
)

func (s *Server) UserGet(ctx context.Context, q *pb.ID) (*pb.User, error) {
	user, err := postgres.GetUser(q.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return user, nil
}

func (s *Server) UserPosts(ctx context.Context, q *pb.IDQuery) (*pb.Posts, error) {
	posts, err := postgres.GetPostsFromUser(q.Id, q.Offset)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return &pb.Posts{Posts: posts}, nil
}

func (s *Server) UserSpaces(ctx context.Context, q *pb.IDQuery) (*pb.Spaces, error) {
	spaces, err := postgres.GetFollowingSpacesForUser(q.Id, q.Offset)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return &pb.Spaces{Spaces: spaces}, nil
}

func (s *Server) UserFollowers(ctx context.Context, q *pb.IDQuery) (*pb.Users, error) {
	followers, err := postgres.GetFollowersForUser(q.Id, q.Offset)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return &pb.Users{Users: followers}, nil
}

func (s *Server) UserFollowing(ctx context.Context, q *pb.IDQuery) (*pb.Following, error) {
	followingUsers, err := postgres.GetFollowingForUser(q.Id, q.Offset)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return &pb.Following{Users: followingUsers}, nil
}
