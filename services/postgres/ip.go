package postgres

import (
	"github.com/rs/zerolog/log"
)

const (
	// TableIP name const
	TableIP = "ip"
	// ColIP const
	ColIP = "ip"
	// ColGeoliteID const
	ColGeoliteID = "geolite_id"
	// ColCellID const
	ColCellID = "cell_id"
)

// GetIPAddressPlace gets the Place corresponding to the IP.
func GetSpaceForIP(ip string) (int64, error) {
	/* We remove this check as we assume all calls have valid IP and the DB trip is worth it.
	if len(ip) == 0 {
		return 0, errors.ErrInvalidArgument
	}*/
	var space int64
	if err := DB.QueryRow(`SELECT space
		FROM ip
		WHERE ip = $1
		LIMIT 1`, ip).Scan(&space); err != nil {
		log.Error().Err(err).Str("ip", ip).Msg("error getting space for ip")
		return 0, err
	}
	return space, nil
}
