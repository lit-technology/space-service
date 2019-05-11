package postgres

import (
	pb "github.com/philip-bui/space-service/protos"
	"github.com/rs/zerolog/log"
)

const (
	TableUserFollower  = "user_follower"
	TableUserFollowing = "user_following"
)

var (
	FollowUserStmt = PrepareStatement(`
		WITH F AS (
			INSERT INTO user_follow(user_id, follower_id)
			VALUES ($1, $2)
			RETURNING user_id, follower_id
		), U AS (
			UPDATE "user" AS U1
			SET followers = COALESCE(followers, 0) + 1
			FROM F
			WHERE U1.id = F.user_id
		)
		UPDATE "user" AS U2
		SET following = COALESCE(following, 0) + 1
		FROM F
		WHERE U2.id = F.follower_id`)
	UnFollowerUserStmt = PrepareStatement(`
		WITH F AS (
			DELETE FROM user_follow
			WHERE user_id = $1
			AND follower_id = $2
			RETURNING user_id, follower_id
		), U1 AS (
			UPDATE "user" AS U1
			SET followers = followers - 1
			FROM F
			WHERE U1.id = F.user_id
		)
		UPDATE "user" AS U2
		SET following = following - 1
		FROM F
		WHERE U2.id = F.follower_id`)
	FollowersForUserStmt = PrepareStatement(`
		SELECT U.id, U.name, U.username, U.photo_url
		FROM user_follow AS F
		INNER JOIN "user" AS U
			ON U.id = F.follower_id
		WHERE F.user_id = $1
		ORDER BY F.follower_id
		LIMIT 20
		OFFSET $2`)
	FollowingForUserStmt = PrepareStatement(`
		SELECT U.id, U.name, U.username, U.photo_url
		FROM user_follow AS F
		INNER JOIN "user" AS U
			ON U.id = F.user_id
		WHERE F.follower_id = $1
		ORDER BY F.user_id
		LIMIT 20
		OFFSET $2`)
	FollowingSpacesForUserStmt = PrepareStatement(`
		SELECT S.id, S.name, S.photo_url
		FROM space_follow AS F
		INNER JOIN space AS S
			ON S.id = F.space_id
		WHERE F.user_id = $1
		ORDER BY F.space_id
		LIMIT 20
		OFFSET $2`)
)

func FollowUser(userID, otherID int64) error {
	if _, err := FollowUserStmt.Exec(userID, otherID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("otherID", otherID).Msg("error following user")
		return err
	}
	return nil
}

func UnFollowUser(userID, otherID int64) error {
	if _, err := UnFollowerUserStmt.Exec(userID, otherID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("otherID", otherID).Msg("error un-following user")
		return err
	}
	return nil
}

func GetFollowersForUser(userID int64, offset int32) ([]*pb.User, error) {
	rows, err := FollowingForUserStmt.Query(userID, offset)
	if err != nil {
		log.Error().Err(err).Int64("userID", userID).Int32("offset", offset).Msg("error getting followers for user")
		return nil, err
	}
	users := make([]*pb.User, 20)
	count := 0
	for rows.Next() {
		u := &pb.User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Username, &u.PhotoURL); err != nil {
			log.Error().Err(err).Int64("userID", userID).Int32("offset", offset).Msg("error scanning followers for user")
			return nil, err
		}
		users = append(users, u)
		count++
	}
	return users[:count], nil
}

func GetFollowingForUser(userID int64, offset int32) ([]*pb.User, error) {
	rows, err := FollowingForUserStmt.Query(userID, offset)
	if err != nil {
		log.Error().Err(err).Int64("userID", userID).Int32("offset", offset).Msg("error getting following for user")
		return nil, err
	}
	users := make([]*pb.User, 20)
	count := 0
	for rows.Next() {
		u := &pb.User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Username, &u.PhotoURL); err != nil {
			log.Error().Err(err).Int64("userID", userID).Int32("offset", offset).Msg("error scanning following for user")
			return nil, err
		}
		users = append(users, u)
		count++
	}
	return users[:count], nil
}

func GetFollowingSpacesForUser(spaceID int64, offset int32) ([]*pb.Space, error) {
	rows, err := FollowingSpacesForUserStmt.Query(spaceID, offset)
	if err != nil {
		log.Error().Err(err).Int64("spaceID", spaceID).Int32("offset", offset).Msg("error getting following spaces for user")
		return nil, err
	}
	spaces := make([]*pb.Space, 20)
	count := 0
	for rows.Next() {
		s := &pb.Space{}
		if err := rows.Scan(&s.Id, &s.Name, &s.PhotoURL); err != nil {
			log.Error().Err(err).Int64("spaceID", spaceID).Int32("offset", offset).Msg("error scanning following spaces for user")
			return nil, err
		}
		spaces = append(spaces, s)
		count++
	}
	return spaces[:count], nil
}
