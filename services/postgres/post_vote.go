package postgres

import (
	"github.com/rs/zerolog/log"
)

var (
	PostUpvoteStmt = PrepareStatement(`
		WITH V AS (
			INSERT INTO post_vote
			VALUES ($1, $2, true)
			RETURNING post_id
		)
		UPDATE post AS P
		SET upvote = upvote + 1
		FROM V
		WHERE P.id = V.post_id`)
	PostUnUpvoteStmt = PrepareStatement(`
		WITH V AS (
			DELETE FROM post_vote
			WHERE user_id = $1
			AND post_id = $2
			AND vote = true
			RETURNING post_id
		)
		UPDATE post AS P
		SET upvote = upvote - 1
		FROM V
		WHERE P.id = V.post_id`)
	PostDownvoteStmt = PrepareStatement(`
		WITH V AS (
			INSERT INTO post_vote
			VALUES ($1, $2, false)
			RETURNING post_id
		)
		UPDATE post AS P
		SET downvote = downvote + 1
		FROM V
		WHERE P.id = V.post_id`)
	PostUnDownvoteStmt = PrepareStatement(`
		WITH V AS (
			DELETE FROM post_vote
			WHERE user_id = $1
			AND post_id = $2
			AND vote = false
			RETURNING post_id
		)
		UPDATE post AS P
		SET downvote = downvote - 1
		FROM V
		WHERE P.id = V.post_id`)
)

func PostUpvote(userID, postID int64) error {
	if _, err := PostUpvoteStmt.Exec(userID, postID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("postID", postID).Msg("error upvoting post")
		return err
	}
	return nil
}

func PostUnUpvote(userID, postID int64) error {
	if _, err := PostUnUpvoteStmt.Exec(userID, postID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("postID", postID).Msg("error un-upvoting post")
		return err
	}
	return nil
}

func PostDownvote(userID, postID int64) error {
	if _, err := PostDownvoteStmt.Exec(userID, postID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("postID", postID).Msg("error downvoting post")
		return err
	}
	return nil
}

func PostUnDownvote(userID, postID int64) error {
	if _, err := PostUnDownvoteStmt.Exec(userID, postID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("postID", postID).Msg("error un-downvoting post")
		return err
	}
	return nil
}
