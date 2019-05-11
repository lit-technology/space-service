package postgres

import (
	"github.com/philip-bui/grpc-errors"
	"github.com/rs/zerolog/log"
)

const (
	// TableGeolite name const
	TableGeolite = "geolite"
)

var (
	GetSpaceIDByCountryStmt = PrepareStatement("SELECT space_id FROM geolite WHERE country_name = $1 AND state_name IS NULL AND city_name IS NULL")
	GetSpaceIDByStateStmt   = PrepareStatement("SELECT space_id FROM geolite WHERE state_name = $1 AND city_name IS NULL AND country_name = $2")
	GetSpaceIDByCityStmt    = PrepareStatement("SELECT space_id FROM geolite WHERE city_name = $1 AND state_name = $2 AND country_name = $3")
)

func GetSpaceIDByCountry(country string) (int64, error) {
	if len(country) == 0 {
		return 0, errors.NewInvalidArgumentError("invalid country")
	}
	var id int64
	err := GetSpaceIDByCountryStmt.QueryRow(country).Scan(&id)
	if err != nil {
		log.Error().Str("country", country).Err(err).Msg("error getting space id by country")
	}
	return id, nil
}

func GetSpaceIDByState(state, country string) (int64, error) {
	if len(state) == 0 || len(country) == 0 {
		return 0, errors.NewInvalidArgumentError("invalid state or country")
	}
	var id int64
	err := GetSpaceIDByStateStmt.QueryRow(state, country).Scan(&id)
	if err != nil {
		log.Error().Str("country", country).Str("state", state).Err(err).Msg("error getting space id by state")
	}
	return id, err
}

func GetSpaceIDByCity(city, state, country string) (int64, error) {
	if len(city) == 0 || len(state) == 0 || len(country) == 0 {
		return 0, errors.NewInvalidArgumentError("invalid city, state or country")
	}
	var id int64
	err := GetSpaceIDByCityStmt.QueryRow(city, state, country).Scan(&id)
	if err != nil {
		log.Error().Str("country", country).Str("state", state).Str("city", city).Err(err).Msg("error getting space id by city")
	}
	return id, err
}
