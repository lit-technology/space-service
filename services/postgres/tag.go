package postgres

import (
	pb "github.com/philip-bui/space-service/protos"
	"github.com/rs/zerolog/log"
)

var (
	InsertTagStmt = PrepareStatement("INSERT INTO tag (tag) VALUES ($1) ON CONFLICT DO NOTHING")
	GetTagStmt    = PrepareStatement("SELECT posts, followers FROM tag WHERE tag = $1")
	SearchTagStmt = PrepareStatement("SELECT tag FROM tag WHERE tag LIKE $1 LIMIT 20")
)

func InsertTag(tag string) error {
	if _, err := InsertTagStmt.Exec(tag); err != nil {
		log.Error().Err(err).Str("tag", tag).Msg("error inserting tag")
		return err
	}
	return nil
}

func GetTag(tag string) (*pb.Tag, error) {
	t := &pb.Tag{Tag: tag}
	if err := GetTagStmt.QueryRow(tag).Scan(&t.Posts, &t.Followers); err != nil {
		log.Error().Err(err).Str("tag", tag).Msg("error getting tag")
		return nil, err
	}
	return t, nil
}

func SearchTag(tag string) ([]*pb.Tag, error) {
	rows, err := SearchTagStmt.Query(tag + "%")
	if err != nil {
		log.Error().Err(err).Str("tag", tag).Msg("error searching for tag")
		return nil, err
	}
	tags := make([]*pb.Tag, 20)
	count := 0
	for rows.Next() {
		t := &pb.Tag{}
		rows.Scan(&t.Tag, &t.Posts, &t.Followers)
		tags = append(tags, t)
		count++
	}
	return tags[:count], nil
}
