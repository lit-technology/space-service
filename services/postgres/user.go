package postgres

import (
	pb "github.com/philip-bui/space-service/protos"
	"github.com/rs/zerolog/log"
)

var (
	InsertUserStmt = PrepareStatement(`
		INSERT INTO "user" (username)
		VALUES ($1)`)
	GetUserStmt = PrepareStatement(`
		SELECT name, username, bio, photo_url, posts, followers, following
		FROM "user"
		WHERE id = $1
		AND deleted < NOW()`)
	SearchUserStmt = PrepareStatement(`
		SELECT id, name, username, photo_url, posts, followers, following
		FROM "user"
		WHERE username LIKE $1
		LIMIT 20`)
)

func InsertUser(username string) error {
	if _, err := InsertUserStmt.Exec(username); err != nil {
		log.Error().Err(err).Str("username", username).Msg("error inserting user")
		return err
	}
	return nil
}

func GetUser(userID int64) (*pb.User, error) {
	u := &pb.User{}
	if err := GetUserStmt.QueryRow(userID).Scan(&u.Name, &u.Username, &u.PhotoURL, &u.Posts, &u.Followers, &u.Following); err != nil {
		log.Error().Err(err).Int64("userID", userID).Msg("error getting user")
		return nil, err
	}
	return u, nil
}

func SearchUser(username string) ([]*pb.User, error) {
	rows, err := SearchUserStmt.Query(username + "%")
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("error searching for user")
		return nil, err
	}
	users := make([]*pb.User, 20)
	count := 0
	for rows.Next() {
		u := &pb.User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Username, &u.PhotoURL, &u.Posts, &u.Followers, &u.Following); err != nil {
			log.Error().Err(err).Msg("error scanning search for user")
			break
		}
		users = append(users, u)
	}
	return users[:count], nil
}
