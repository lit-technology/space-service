package controllers

import (
	"context"
	"database/sql"
	"strings"

	"github.com/philip-bui/grpc-errors"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/philip-bui/space-service/services/auth"
	"github.com/philip-bui/space-service/services/postgres"
)

func (s *Server) SpaceNew(ctx context.Context, req *pb.SpaceNewRequest) (*pb.SpaceNewResponse, error) {
	if err := SpaceNewValidate(req); err != nil {
		return nil, err
	}
	_, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if len(req.Photo) != 0 {

	}
	//postgres.InsertSpaceWithPlace(userID, req.Name)
	return nil, nil
}

func SpaceNewValidate(req *pb.SpaceNewRequest) error {
	name := strings.TrimSpace(req.Name)
	if len(name) <= 6 || len(name) > 180 {
		return errors.NewInvalidArgumentError("invalid name")
	}
	if req.Place == nil || len(req.Place.Country) == 0 {
		return errors.NewInvalidArgumentError("invalid place")
	}
	SpaceNewGenderValidate(req.Gender)
	return nil
}

func SpaceNewGenderValidate(gender pb.Gender) interface{} {
	switch gender {
	case pb.Gender_NONE:
		return postgres.NullInt64
	case pb.Gender_FEMALE:
		return false
	case pb.Gender_MALE:
		return true
	}
	return postgres.NullInt64
}

func (s *Server) SpaceGet(ctx context.Context, q *pb.ID) (*pb.Space, error) {
	space, err := postgres.GetSpaceWithClaims(q.Id, auth.GetClaimsFromContextUnsafe(ctx))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return space, nil
}

func (s *Server) SpacePosts(ctx context.Context, q *pb.IDQuery) (*pb.Posts, error) {
	posts, err := postgres.GetPostsForSpace(q.Id, q.Offset)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return &pb.Posts{Posts: posts}, nil
}

func (s *Server) SpaceSpaces(ctx context.Context, q *pb.IDQuery) (*pb.Spaces, error) {
	// TODO:
	return nil, nil
}
