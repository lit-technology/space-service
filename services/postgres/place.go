package postgres

import (
	"database/sql"

	"github.com/philip-bui/grpc-errors"
	pb "github.com/philip-bui/space-service/protos"
	"github.com/rs/zerolog/log"
)

const (
	// ColStreetName const
	ColStreetName = "street_name"
	// ColPlaceName const
	ColPlaceName = "place_name"
)

func InsertOrGetSpaceIDByPlace(p *pb.Place) (int64, error) {
	switch p.Type {
	case pb.Place_COUNTRY:
		return GetSpaceIDByCountry(p.Country)
	case pb.Place_STATE:
		return GetSpaceIDByState(p.State, p.Country)
	case pb.Place_CITY:
		return GetSpaceIDByCity(p.City, p.State, p.Country)
	case pb.Place_ADDRESS:
		id, err := GetSpaceIDByAddress(p.Street, p.City, p.State, p.Country)
		if err != nil && err == sql.ErrNoRows {
			/*cityID, err := GetSpaceIDByCity(p.City, p.State, p.Country)
			if err != nil {
				return 0, err
			}*/
		}
		return id, err
	case pb.Place_PLACE:
		id, err := GetSpaceIDByPlace(p.Name, p.Street, p.City, p.State, p.Country)
		if err != nil && err == sql.ErrNoRows {

		}
		return id, err
	}
	return 0, errors.ErrInvalidArgument
}

func GetSpaceIDByAddress(street, city, state, country string) (int64, error) {
	if len(street) == 0 || len(city) == 0 || len(state) == 0 || len(country) == 0 {
		return 0, errors.NewInvalidArgumentError("invalid address, city, state or country")
	}
	var id int64
	err := DB.QueryRow("SELECT space FROM place WHERE street_name = $1 AND city_name = $2 AND state_name = $3 AND country_name = $4", street, city, state, country).Scan(&id)
	if err != nil {
		log.Error().Str("country", country).Str("state", state).Str("city", city).Str("street", street).Err(err).Msg("error getting space id by address")
	}
	return id, err
}

func InsertAddress(street, city, state, country string) (int64, error) {
	if len(street) == 0 || len(city) == 0 || len(state) == 0 || len(country) == 0 {
		return 0, errors.NewInvalidArgumentError("invalid address, city, state or country")
	}
	var id int64
	err := DB.QueryRow("INSERT INTO place(space, street, country, state, city) VALUES($1, $2, $3, $4) RETURNING id", street, city, state, country).Scan(&id)
	if err != nil {
		log.Error().Str("country", country).Str("state", state).Str("city", city).Str("street", street).Err(err).Msg("error getting space id by address")
	}
	return id, err
}

func GetSpaceIDByPlace(place, street, city, state, country string) (int64, error) {
	if len(street) == 0 || len(city) == 0 || len(state) == 0 || len(country) == 0 {
		return 0, errors.NewInvalidArgumentError("invalid place, address, city, state or country")
	}
	var id int64
	err := DB.QueryRow("SELECT space FROM place WHERE place = $1 AND street_name = $2 AND city_name = $3 AND state_name = $4 AND country_name = $5", place, street, city, state, country).Scan(&id)
	if err != nil {
		log.Error().Str("country", country).Str("state", state).Str("city", city).Str("street", street).Str("place", place).Err(err).Msg("error getting space id by place")
	}
	return id, err
}
