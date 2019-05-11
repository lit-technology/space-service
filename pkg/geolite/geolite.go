package geolite

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/lib/pq"
	"github.com/philip-bui/space-service/pkg/csv"
	"github.com/philip-bui/space-service/services/postgres"
	"github.com/rs/zerolog/log"
)

// GeoliteCityRow is a struct matching a GeoLite2 CSV City row.
type GeoliteCityRow struct {
	GeonameID           int    `csv:"geoname_id"`
	ContinentCode       string `csv:"continent_code"`
	ContinentName       string `csv:"continent_name"`
	CountryISOCode      string `csv:"country_iso_code"`
	CountryName         string `csv:"country_name"`
	Subdivision1ISOCode string `csv:"subdivision_1_iso_code"`
	Subdivision1Name    string `csv:"subdivision_1_name"`
	Subdivision2ISOCode string `csv:"subdivision_2_iso_code"`
	Subdivision2Name    string `csv:"subdivision_2_name"`
	CityName            string `csv:"city_name"`
	MetroCode           string `csv:"metro_code"`
	TimeZone            string `csv:"time_zone"`
	IsInEuropeanUnion   bool   `csv:"is_in_european_union"`
}

var (
	countries  = make(map[string]int64, 300)
	countryIDs = int64(0)
	states     = make(map[string]int64, 3000)
	stateIDs   = make(map[int64]int64, 3000)
	cities     = make(map[string]int64, 18000)
	cityIDs    = make(map[string]int64, 18000)
)

// ReadCities reads a Geolite2 CSV City file, parses into Continent, Country,
// State, City into PostgreSQL and batches matching GeoliteID likewise.
func ReadCities(fName string) {
	f, err := os.Open(fName)
	defer f.Close()
	if err != nil {
		log.Fatal().Err(err).Str("fName", fName).Msg("error opening file")
	}
	row := &GeoliteCityRow{}
	u := csv.NewCsvUnmarshallerFromFile(f, row)
	batch, err := BeginGeoliteBatch()
	if err != nil {
		log.Fatal().Err(err).Msg("error creating geolite batch")
	}
	for err = u.UnmarshalToStruct(row); err == nil; err = u.UnmarshalToStruct(row) {
		ID := int64(0)
		name := row.CountryName
		// Country
		countryID, ok := countries[row.CountryISOCode]
		if !ok {
			countryIDs++
			countries[row.CountryISOCode] = countryIDs
			countryID = countryIDs
		}
		ID = ID + (countryID << 56)
		countryIDString := strconv.FormatInt(countryID, 10)
		var stateCode interface{} = sql.NullString{}
		var stateName interface{} = sql.NullString{}
		var areaCode interface{} = sql.NullString{}
		var areaName interface{} = sql.NullString{}
		var cityName interface{} = sql.NullString{}
		var stateID int64 = 0
		var cityID int64 = 0
		// State
		if len(row.Subdivision1Name) != 0 || len(row.Subdivision2Name) != 0 {
			if len(row.Subdivision2Name) != 0 {
				areaCode = row.Subdivision2ISOCode
				areaName = row.Subdivision2Name
				name = row.Subdivision2Name
			}
			if len(row.Subdivision1Name) != 0 {
				stateCode = row.Subdivision1ISOCode
				stateName = row.Subdivision1Name
				name = row.Subdivision1Name
			}
			if len(row.Subdivision1Name) != 0 && len(row.Subdivision2Name) != 0 {
				name = row.Subdivision1Name + " " + row.Subdivision2Name
			}
			stateKey := countryIDString + row.Subdivision1Name + row.Subdivision2Name
			stateID, ok = states[stateKey]
			if !ok {
				stateIDs[countryID]++
				states[stateKey] = stateIDs[countryID]
				stateID = stateIDs[countryID]
			}
			ID = ID + (stateID << 48)
			// City
			if len(row.CityName) != 0 {
				cityName = row.CityName
				name = row.CityName
				stateIDString := strconv.FormatInt(stateID, 10)
				cityKey := stateIDString + row.CityName
				cityID, ok = cities[cityKey]
				if !ok {
					cityIDs[stateKey]++
					cities[cityKey] = cityIDs[stateKey]
					cityID = cityIDs[stateKey]
				}
				ID = ID + (cityID << 37)
			}
		} else if len(row.CityName) != 0 {
			cityName = row.CityName
			name = row.CityName
			stateIDString := strconv.FormatInt(stateID, 10)
			cityKey := stateIDString + row.CityName
			cityID, ok = cities[cityKey]
			if !ok {
				cityIDs[countryIDString]++
				cities[cityKey] = cityIDs[countryIDString]
				cityID = cityIDs[countryIDString]
			}
			ID = ID + (cityID << 37)
		}
		spaceID, err := postgres.InsertSpace(ID, name, countryID, stateID, cityID)
		if err != nil {
			log.Warn().Int64("country", countryID).Int64("state", stateID).Int64("city", cityID).Msgf("error creating space %v", row)
			spaceID = ID
		}
		batch.AddRow(row.GeonameID, spaceID, row.ContinentCode, row.ContinentName, row.CountryISOCode, row.CountryName, stateCode, stateName, areaCode, areaName, cityName)
	}
	if err := batch.ExecAndCommit(); err != nil {
		log.Fatal().Err(err).Msg("error commiting geolite batch")
	}
}

// GeoliteBatch is a transaction handling start, batching of Geolites and commit of the transaction.
type GeoliteBatch struct {
	*sql.Tx
	*sql.Stmt
}

// BeginGeoliteBatch begins a Geolite Batch Transaction.
func BeginGeoliteBatch() (*GeoliteBatch, error) {
	txn, err := postgres.DB.Begin()
	if err != nil {
		log.Error().Err(err).Msg("error beginning transaction")
		return nil, err
	}
	stmt, err := txn.Prepare(pq.CopyIn(postgres.TableGeolite, postgres.ColID, postgres.ColSpaceID,
		postgres.ColContinentCode, postgres.ColContinentName, postgres.ColCountryCode, postgres.ColCountryName,
		postgres.ColStateCode, postgres.ColStateName, postgres.ColAreaCode, postgres.ColAreaName, postgres.ColCityName))
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return nil, err
	}
	return &GeoliteBatch{
		txn,
		stmt,
	}, nil
}

// AddRow adds a Geolite row.
func (g *GeoliteBatch) AddRow(id int, spaceID int64, continentCode, continentName, countryCode, countryName string, stateCode, stateName, areaCode, areaName, cityName interface{}) error {
	if _, err := g.Stmt.Exec(id, spaceID, continentCode, continentName, countryCode, countryName, stateCode, stateName, areaCode, areaName, cityName); err != nil {
		log.Error().Err(err).Int("id", id).Int64("spaceID", spaceID).Msg("error adding row")
		return err
	}
	return nil
}

// ExecAndCommit executes and commits the batch transaction.
func (g *GeoliteBatch) ExecAndCommit() error {
	if _, err := g.Stmt.Exec(); err != nil {
		log.Error().Err(err).Msg("error executing statement")
		return err
	}
	if err := g.Close(); err != nil {
		log.Error().Err(err).Msg("error closing statement")
		return err
	}
	if err := g.Commit(); err != nil {
		log.Error().Err(err).Msg("error commiting transaction")
		return err
	}
	return nil
}
