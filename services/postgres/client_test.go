package postgres

import (
	"database/sql"
	"testing"

	"github.com/philip-bui/space-service/config"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

const (
	SpaceID = 7777
	PostID  = 77
)

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}

type PostgresSuite struct {
	suite.Suite
}

func (s *PostgresSuite) SetupSuite() {
	DB.Exec("CREATE DATABASE test") // Don't worry about errors here.
	_, err := sql.Open("postgres", "host="+config.PostgresHost+
		" port="+config.PostgresPort+
		" user="+config.PostgresUser+
		" password="+config.PostgresPass+
		" dbname=Test"+
		" sslmode=disable")
	if err != nil {
		s.FailNow(err.Error())
	}
	//DB = db
	zerolog.SetGlobalLevel(zerolog.FatalLevel)
}

func (s *PostgresSuite) BeforeTest() {
	DB.Exec("DELETE * FROM post")
}
