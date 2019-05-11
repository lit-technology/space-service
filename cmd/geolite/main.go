package main

import (
	"github.com/philip-bui/space-service/pkg/geolite"
	//	"github.com/philip-bui/space-service/services/postgres"
	//	"github.com/rs/zerolog/log"
)

func main() {
	geolite.ReadCities("GeoLite2-City-Locations-en.csv")
	geolite.ReadIPv4("GeoLite2-City-Blocks-IPv4.csv")
	geolite.ReadIPv6("GeoLite2-City-Blocks-IPv6.csv")
}
