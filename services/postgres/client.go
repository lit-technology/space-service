package postgres

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	"github.com/philip-bui/space-service/config"
	"github.com/rs/zerolog/log"
)

var (
	// DB is the exported PostgreSQL Client.
	DB         = initDB()
	NullInt64  = sql.NullInt64{}
	NullBool   = sql.NullBool{}
	NullString = sql.NullString{}
)

const (
	EmptyValue = " "
	// ColID const
	ColID = "id"
	// ColSpaceID const
	ColSpaceID = "space_id"
	// ColUserId
	ColUserID = "user_id"
	// ColContinentCode const
	ColContinentCode = "continent_code"
	// ColContinentName const
	ColContinentName = "continent_name"
	// ColCountryCode const
	ColCountryCode = "country_code"
	// ColCountryName const
	ColCountryName = "country_name"
	// ColStateCode const
	ColStateCode = "state_code"
	// ColStateName const
	ColStateName = "state_name"
	// ColAreaCode const
	ColAreaCode = "area_code"
	// ColAreaName const
	ColAreaName = "area_name"
	// ColCityName const
	ColCityName = "city_name"
)

func initDB() *sql.DB {
	db, err := sql.Open("postgres", "host="+config.PostgresHost+
		" port="+config.PostgresPort+
		" user="+config.PostgresUser+
		" password="+config.PostgresPass+
		" dbname="+config.PostgresDB+
		" sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to postgres")
	}
	db.SetMaxIdleConns(0)
	return db
}

func NullableEmptyStringArray(a []string, sep string) interface{} {
	if a == nil || len(a) == 0 {
		return NullString
	}
	s := strings.Join(a, sep)
	return NullableEmptyString(s)
}

func NullableEmptyString(s string) interface{} {
	if s == "" {
		return NullString
	}
	return s
}

func NullableZero(i int64) interface{} {
	if i == 0 {
		return NullInt64
	}
	return i
}

func NullableBool(b bool) interface{} {
	if !b {
		return NullBool
	}
	return b
}

func PrepareStatement(query string) *sql.Stmt {
	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Fatal().Err(err).Str("query", strings.Replace(query, "\t", "", 100)).Msg("error preparing statement")
	}
	return stmt
}
