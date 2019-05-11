package main

import (
	"net"
	"os"
	"time"

	"github.com/philip-bui/space-service/config"
	"github.com/philip-bui/space-service/controllers"
	pb "github.com/philip-bui/space-service/protos"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC1123
	LOG := os.Getenv("LOG")
	switch LOG {
	case "ERROR":
		log.Logger = log.Level(zerolog.ErrorLevel)
	case "WARN":
		log.Logger = log.Level(zerolog.WarnLevel)
	case "INFO":
		log.Logger = log.Level(zerolog.InfoLevel)
	default:
		log.Logger = log.Level(zerolog.DebugLevel)
	}
	if log.Debug().Enabled() {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func main() {
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("listen to tcp")
	} else {
		log.Info().Msg("listen to tcp")
	}
	s := grpc.NewServer()
	reflection.Register(s)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	pb.RegisterSparkServer(s, controllers.NewServer())
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("serve grpc")
	}
}
