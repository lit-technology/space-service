package postgres

import (
	"database/sql"
	"strconv"

	//	"github.com/lib/pq"
	"github.com/philip-bui/space-service/pkg/bits"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/rs/zerolog/log"
)

// SpaceType corresponds to a Country, State, City, Address, Place or Group.
type SpaceType int16

const (
	// SpaceCountry corresponds to a country.
	SpaceCountry SpaceType = iota
	// SpaceState corresponds to a state.
	SpaceState
	// SpaceCity corresponds to a city.
	SpaceCity
	// SpaceAddress corresponds to a address.
	SpaceAddress
	// SpaceCity corresponds to a place.
	SpacePlace
	// SpaceGroup corresponds to a group.
	SpaceGroup
)

const (
	SpaceIDCityBitshiftLeft = 64 - 8 - 8 - 11
)

const (
	CountryMask = 0x8
	CountryBits = 8
	StateMask   = CountryMask
	StateBits   = CountryBits
	CityMask    = 0xB
	CityBits    = 11
	PlaceMask   = 0xFFFFFFFFFF
	PlaceBits   = 40
)

var (
	GetSpaceStmt = PrepareStatement(`
		SELECT name
		FROM space
		WHERE id = $1
		AND deleted < NOW()`)
	GetSpacesForSpaceStmt = PrepareStatement(`
		SELECT name, gender, age_min, age_max, followers
		FROM space
		WHERE id BETWEEN $1 AND $2
		ORDER BY followers
		LIMIT 20
		OFFSET $3`)
)

// GetSpaceIDByName gets the countryID using its name.
func GetSpaceIDByName(name string) (int64, error) {
	var id int64
	err := DB.QueryRow("SELECT id FROM space WHERE name = $1", name).Scan(&id)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("error getting space")
	}
	return id, err
}

func InsertSpaceForUser(userID, name, photoURL, password, bio, country, state, city string, gender interface{}, ageMin, ageMax int32) (int64, error) {

	return 0, nil
}

// InsertSpace inserts a country with continentID, code and name, returning the countryID.
func InsertSpace(ID int64, name string, countryID, stateID, cityID int64) (int64, error) {
	var id int64
	var spaceType SpaceType
	if cityID != 0 {
		spaceType = SpaceCity
	} else if stateID != 0 {
		spaceType = SpaceState
	} else {
		spaceType = SpaceCountry
	}
	err := DB.QueryRow("INSERT INTO space(id, name, country, state, city, type) VALUES($1, $2, $3, $4, $5, $6) RETURNING id", ID, name, countryID, NullableZero(stateID), NullableZero(cityID), spaceType).Scan(&id)
	if err != nil {
		log.Error().Err(err).Int64("ID", ID).Str("name", name).Int64("countryID", countryID).Int64("stateID", stateID).Int64("cityID", cityID).Msg("error inserting space")
	}
	return id, err
}

func GetSpaceWithClaims(spaceID int64, claims map[string]interface{}) (*pb.Space, error) {
	s := &pb.Space{}

	return s, nil
}

func GetPlacesForSpaceWithClaims(spaceID int64, offset int32, claims map[string]interface{}) ([]*pb.Space, error) {
	rows, err := GetPlacesForSpace(spaceID, offset)
	if err != nil {
		return nil, err
	}
	return ScanSpacesWithClaims(rows, claims)
}

func GetChildrenForSpace(spaceID int64, offset int32) ([]*pb.Space, error) {
	if !bits.IsUnsetBitsFromRight(CityMask, spaceID, PlaceBits) {

	} else if !bits.IsUnsetBitsFromRight(StateMask, spaceID, PlaceBits+CityBits) {

	} else {

	}
	return nil, nil
}

func ScanSpacesWithClaims(rows *sql.Rows, claims map[string]interface{}) ([]*pb.Space, error) {
	spaces := make([]*pb.Space, 20)
	count := 0
	for rows.Next() {
		s := &pb.Space{}
		password := sql.NullString{}
		gender := sql.NullBool{}
		ageMin := sql.NullInt64{}
		ageMax := sql.NullInt64{}
		if err := rows.Scan(&s.Id, &s.Name, &password, &gender, &ageMin, &ageMax, &s.Posts); err != nil {
			log.Error().Err(err).Msg("error scanning spaces")
		}
		if _, ok := claims["g"].(bool); ok {

		}
		spaces[count] = s
		count++
	}
	return spaces[:count], nil
}

func GetPlacesForSpace(spaceID int64, offset int32) (*sql.Rows, error) {
	cityID, nextCityID := bits.RangeInt64(PlaceBits, spaceID)
	rows, err := GetSpacesForSpaceStmt.Query(spaceID, cityID, nextCityID, offset)
	if err != nil {
		log.Error().Err(err).Int64("spaceID", spaceID).
			Str("start", strconv.FormatInt(cityID, 2)).
			Str("end", strconv.FormatInt(nextCityID, 2)).
			Int32("offset", offset).Msg("error getting places for space")
		return nil, err
	}
	return rows, err
}

func GetCityAndPlacesForSpace(spaceID int64, offset int32) (*sql.Rows, error) {
	stateID, nextStateID := bits.RangeInt64(PlaceBits+CityBits, spaceID)
	rows, err := GetSpacesForSpaceStmt.Query(spaceID, stateID, nextStateID, offset)
	if err != nil {
		log.Error().Err(err).Int64("spaceID", spaceID).
			Str("start", strconv.FormatInt(stateID, 2)).
			Str("end", strconv.FormatInt(nextStateID, 2)).
			Int32("offset", offset).Msg("error getting cities and places for space")
		return nil, err
	}
	return rows, err
}

func GetStateCityAndPlacesForCountry(spaceID int64, offset int32) (*sql.Rows, error) {
	countryID, nextCountryID := bits.RangeInt64(PlaceBits+CityBits+StateBits, spaceID)
	rows, err := GetSpacesForSpaceStmt.Query(spaceID, countryID, nextCountryID, offset)
	if err != nil {
		log.Error().Err(err).Int64("spaceID", spaceID).
			Str("start", strconv.FormatInt(countryID, 2)).
			Str("end", strconv.FormatInt(nextCountryID, 2)).
			Int32("offset", offset).Msg("error getting states, cities and places for space")
		return nil, err
	}
	return rows, err
}
