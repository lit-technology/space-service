package auth

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/peer"
)

func GetIpFromContext(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	} else {
		log.Error().Msg("error retrieving peer from context")
		return ""
	}
}
