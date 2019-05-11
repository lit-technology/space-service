package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/rs/zerolog/log"
)

const (
	TableCommentVote = "comment_vote"
)

var (
	InsertCommentStmt = PrepareStatement(`
		INSERT INTO comment (post_id, user_id, reply_id, text)
		VALUES ($1, $2, $3, $4)`)
	PostCommentsStmt = PrepareStatement(`
		SELECT id, user_id, reply_id, text, upvote, downvote, created
		FROM comment
		WHERE post_id = $1
		ORDER BY upvote
		LIMIT 20
		OFFSET $2`)
	GetCommentsStmt = PrepareStatement(`
		SELECT id, user_id, reply_id, text, upvote, downvote, created
		FROM comment
		WHERE post_id = $1
		AND id IN ($2)`)
	CommentUpvoteStmt = PrepareStatement(`
		WITH V AS (
			INSERT INTO comment_vote
			VALUES ($1, $2, true)
			RETURNING comment_id
		)
		UPDATE comment AS C
		SET upvote = upvote + 1
		FROM V
		WHERE C.id = V.comment_id`)
	CommentUnUpvoteStmt = PrepareStatement(`
		WITH V AS (
			DELETE FROM comment_vote
			WHERE user_id = $1
			AND comment_id = $2
			AND vote = true
			RETURNING comment_id
		)
		UPDATE comment AS C
		SET upvote = upvote - 1
		FROM V
		WHERE C.id = V.comment_id`)
	CommentDownvoteStmt = PrepareStatement(`
		WITH V AS (
			INSERT INTO comment_vote
			VALUES ($1, $2, false)
			RETURNING comment_id
		)
		UPDATE comment AS C
		SET downvote = downvote + 1
		FROM V
		WHERE C.id = V.comment_id`)
	CommentUnDownvoteStmt = PrepareStatement(`
		WITH V AS (
			DELETE FROM comment_vote
			WHERE user_id = $1
			AND comment_id = $2
			AND vote = false
			RETURNING comment_id
		)
		UPDATE comment AS C
		SET downvote = downvote - 1
		FROM V
		WHERE C.id = V.comment_id`)
)

func InsertComment(postID, userID, replyID int64, text string) (int64, error) {
	var commentID int64
	if err := InsertCommentStmt.QueryRow(postID, userID, NullableZero(replyID), text).Scan(&commentID); err != nil {
		log.Error().Err(err).Int64("postID", postID).Int64("userID", userID).Int64("replyID", replyID).Str("text", text).Msg("error inserting comment")
		return 0, err
	}
	return commentID, nil
}

func ScanComments(rows *sql.Rows, limit int) ([]*pb.Comment, error) {
	comments := make([]*pb.Comment, limit)
	count := 0
	for rows.Next() {
		c := &pb.Comment{}
		replyID := sql.NullString{}
		if err := rows.Scan(&c.Id, &c.UserID, &replyID, &c.Text, &c.Upvote, &c.Downvote, &c.Created); err != nil {
			log.Error().Err(err).Msg("error scanning comments")
			return nil, err
		}
		count++
	}
	return comments[:count], nil
}

func GetCommentsForPost(postID int64, offset int32) ([]*pb.Comment, error) {
	rows, err := PostCommentsStmt.Query(postID, offset)
	if err != nil {
		log.Error().Err(err).Int64("postID", postID).Msg("error getting post comments")
		return nil, err
	}
	return ScanComments(rows, 20)
}

func GetComments(postID int64, commentIDs []int64) ([]*pb.Comment, error) {
	rows, err := GetCommentsStmt.Query(postID, pq.Array(commentIDs))
	if err != nil {
		log.Error().Err(err).Int64("postID", postID).Msg("error getting comments")
		return nil, err
	}
	return ScanComments(rows, len(commentIDs))
}

func CommentUpvote(userID, commentID int64) error {
	if _, err := CommentUpvoteStmt.Exec(userID, commentID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("commentID", commentID).Msg("error upvoting comment")
		return err
	}
	return nil
}

func CommentUnUpvote(userID, commentID int64) error {
	if _, err := CommentUnUpvoteStmt.Exec(userID, commentID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("commentID", commentID).Msg("error un-upvoting comment")
		return err
	}
	return nil
}

func CommentDownvote(userID, commentID int64) error {
	if _, err := CommentDownvoteStmt.Exec(userID, commentID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("commentID", commentID).Msg("error downvoting comment")
		return err
	}
	return nil
}

func CommentUnDownvote(userID, commentID int64) error {
	if _, err := CommentUnDownvoteStmt.Exec(userID, commentID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Int64("commentID", commentID).Msg("error un-downvoting comment")
		return err
	}
	return nil
}
