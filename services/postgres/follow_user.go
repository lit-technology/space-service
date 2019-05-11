package postgres

import (
	"github.com/rs/zerolog/log"
)

var (
	FollowSpaceStmt = PrepareStatement(`
		INSERT INTO space_follow(space_id, user_id)
		VALUES ($1, $2)`)
	UnFollowSpaceStmt = PrepareStatement(`
		DELETE FROM space_follow
		WHERE space_id = $1
		AND user_id = $2`)
	FollowTagStmt = PrepareStatement(`
		INSERT INTO tag_follow(tag, user_id)
		VALUES ($1, $2)`)
	UnFollowTagStmt = PrepareStatement(`
		DELETE FROM tag_follow
		WHERE tag = $1
		AND user_id = $2`)
)

func FollowSpace(userID, spaceID int64) error {
	if _, err := FollowSpaceStmt.Exec(spaceID, userID); err != nil {
		log.Error().Err(err).Msg("error following space")
		return err
	}
	return nil
}

func UnFollowSpace(userID, spaceID int64) error {
	if _, err := UnFollowSpaceStmt.Exec(spaceID, userID); err != nil {
		log.Error().Err(err).Msg("error un-following space")
		return err
	}
	return nil
}

func FollowTag(userID int64, tag string) error {
	if _, err := FollowTagStmt.Exec(tag, userID); err != nil {
		log.Error().Err(err).Msg("error following tag")
		return err
	}
	return nil
}

func UnFollowTag(userID int64, tag string) error {
	if _, err := UnFollowTagStmt.Exec(tag, userID); err != nil {
		log.Error().Err(err).Msg("error un-following tag")
		return err
	}
	return nil
}
