package postgres

import (
	"database/sql"
	"strings"
	"time"

	"github.com/lib/pq"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/rs/zerolog/log"
)

const (
	TablePost      = "post"
	TablePostSpace = "post_space"
	TablePostUser  = "post_user"
	ColPostID      = "post_id"
	ColURLs        = "urls"
	ColName        = "name"
	ColNameHidden  = "name_hidden"
	ColEventStart  = "event_start"
	ColEventEnd    = "event_end"
	ColEventPrice  = "event_price"
	ColDeleted     = "deleted"
)

var (
	GetPostStmt           = PrepareStatement("SELECT user_id, name, urls, name_hidden, space_id, upvote, downvote, rating, start_date, end_date, price, created, deleted FROM post WHERE id = $1")
	InsertPostForUserStmt = PrepareStatement(`
		WITH P AS (
			INSERT INTO post(user_id, name,
				urls, name_hidden,
				space_id, rating,
				start_date, end_date,
				price, deleted)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	   		RETURNING id
	   ), PS AS (
	   		INSERT INTO post_space(space_id, post_id)
  			SELECT $5, id
			FROM P
  		)
		INSERT INTO post_user(user_id, post_id)
		SELECT $1, id
		FROM P`)
	GetPostsFromUserStmt = PrepareStatement(`
		SELECT P.id, P.user_id, P.name, P.urls, P.name_hidden, P.space_id, P.upvote, P.downvote, P.rating, P.start_date, P.end_date, P.price, P.created, P.deleted
		FROM post_user AS U
		INNER JOIN post as P
			ON U.post_id = P.id
		WHERE U.user_id = $1
		ORDER BY U.post_id DESC
		LIMIT 20
		OFFSET $2`)
	GetPostsForSpaceStmt = PrepareStatement(`
		SELECT P.id, P.user_id, P.name, P.urls, P.name_hidden, P.space_id, P.upvote, P.downvote, P.rating, P.start_date, P.end_date, P.price, P.created, P.deleted
		FROM post_space AS S
		INNER JOIN post as P
			ON S.post_id = P.id
		WHERE S.space_id= $1
		ORDER BY S.post_id DESC
		LIMIT 20
		OFFSET $2`)
	GetPostsForTagStmt = PrepareStatement(`
		SELECT P.id, P.user_id, P.name, P.urls, P.name_hidden, P.space_id, P.upvote, P.downvote, P.rating, P.start_date, P.end_date, P.price, P.created, P.deleted
		FROM post_tag AS T
		INNER JOIN post as P
			ON T.post_id = P.id
		WHERE T.tag = $1
		ORDER BY T.post_id DESC
		LIMIT 20
		OFFSET $2`)
)

func InsertPostForUser(userID int64, name string, urls []string, nameHidden bool, spaceID, rating, eventStart, eventEnd, eventPrice, expiry int64) error {
	if _, err := InsertPostForUserStmt.Exec(userID, name, strings.Join(urls, ","), NullableBool(nameHidden), spaceID,
		NullableZero(rating), NullableZero(eventStart), NullableZero(eventEnd), NullableZero(eventPrice), NullableZero(expiry)); err != nil {
		log.Error().Err(err).Int64("userID", userID).Msg("error inserting post for user")
		return err
	}
	return nil
}

func GetPost(postID int64) (*pb.Post, error) {
	p := &pb.Post{}
	URLs := sql.NullString{}
	nameHidden := sql.NullBool{}
	rating := sql.NullInt64{}
	startDate := pq.NullTime{}
	endDate := pq.NullTime{}
	price := sql.NullInt64{}
	created := time.Time{}
	deleted := sql.NullInt64{}
	if err := GetPostStmt.QueryRow(postID).Scan(&p.UserID, &p.Name, &URLs, &nameHidden, &p.SpaceID, &p.Upvotes, &p.Downvotes, &rating, &startDate, &endDate, &price, &created, &deleted); err != nil {
		log.Error().Err(err).Int64("postID", postID).Msg("error getting post")
		return nil, err
	}
	if nameHidden.Valid {
		p.NameHidden = nameHidden.Bool
	}
	if URLs.Valid {
		p.URLs = URLs.String
	}
	if rating.Valid {
		p.Rating = float32(rating.Int64)
	}
	if startDate.Valid {
		p.StartDate = startDate.Time.Unix()
	}
	if endDate.Valid {
		p.EndDate = endDate.Time.Unix()
	}
	p.Created = created.Unix()
	if deleted.Valid {
		p.Deleted = deleted.Int64
	}
	return p, nil
}

func ScanPosts(rows *sql.Rows) ([]*pb.Post, error) {
	posts := make([]*pb.Post, 20)
	count := 0
	for rows.Next() {
		p := &pb.Post{}
		URLs := sql.NullString{}
		nameHidden := sql.NullBool{}
		rating := sql.NullInt64{}
		startDate := pq.NullTime{}
		endDate := pq.NullTime{}
		price := sql.NullInt64{}
		created := time.Time{}
		deleted := sql.NullInt64{}
		if err := rows.Scan(&p.Id, &p.UserID, &p.Name, &URLs, &nameHidden, &p.SpaceID, &p.Upvotes, &p.Downvotes, &rating, &startDate, &endDate, &price, &created, &deleted); err != nil {
			log.Error().Err(err).Msg("error scanning posts")
			return nil, err
		}
		if nameHidden.Valid {
			p.NameHidden = nameHidden.Bool
		}
		if URLs.Valid {
			p.URLs = URLs.String
		}
		if rating.Valid {
			p.Rating = float32(rating.Int64)
		}
		if startDate.Valid {
			p.StartDate = startDate.Time.Unix()
		}
		if endDate.Valid {
			p.EndDate = endDate.Time.Unix()
		}
		p.Created = created.Unix()
		if deleted.Valid {
			p.Deleted = deleted.Int64
		}
		posts[count] = p
		count++
	}
	return posts[:count], nil
}

func GetPostsFromUser(userID int64, offset int32) ([]*pb.Post, error) {
	rows, err := GetPostsFromUserStmt.Query(userID, offset)
	if err != nil {
		log.Error().Err(err).Int64("userID", userID).Int32("offset", offset).Msg("error getting post from user")
		return nil, err
	}
	return ScanPosts(rows)
}

func GetPostsForSpace(spaceID int64, offset int32) ([]*pb.Post, error) {
	rows, err := GetPostsForSpaceStmt.Query(spaceID, offset)
	if err != nil {
		log.Error().Err(err).Int64("spaceID", spaceID).Int32("offset", offset).Msg("error getting posts for space")
		return nil, err
	}
	return ScanPosts(rows)
}

func GetPostsForTag(tag string, offset int32) ([]*pb.Post, error) {
	rows, err := GetPostsForTagStmt.Query(tag, offset)
	if err != nil {
		log.Error().Err(err).Str("tag", tag).Int32("offset", offset).Msg("error getting posts for tag")
		return nil, err
	}
	return ScanPosts(rows)
}
