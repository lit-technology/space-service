package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/philip-bui/space-service/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
)

const (
	Authorization = "authorization"
	base          = 32
	sub           = "sub"
)

var (
	ErrNoMetadata     = errors.New("no metadata in context")
	ErrNoAuthMetadata = errors.New("no authorization metadata")
	ErrInvalidToken   = errors.New("invalid token")
)

func GetAuthorizationFromContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if authorization, ok := md[Authorization]; ok && len(authorization) > 0 {
			return strings.Join(authorization, ""), nil
		}
		return "", ErrNoAuthMetadata
	}
	return "", ErrNoMetadata
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	JWT, err := GetAuthorizationFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return GetUserIDFromJWT(JWT)
}

func GetClaimsFromContext(ctx context.Context) (map[string]interface{}, error) {
	JWT, err := GetAuthorizationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return GetClaimsFromJWT(JWT)
}

func GetClaimsFromContextUnsafe(ctx context.Context) map[string]interface{} {
	JWT, err := GetAuthorizationFromContext(ctx)
	if err != nil {
		return make(map[string]interface{}, 0)
	}
	claims, err := GetClaimsFromJWT(JWT)
	if err != nil {
		return make(map[string]interface{}, 0)
	}
	return claims
}

func GetClaimsFromJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		log.Error().Err(err).Str("token", tokenString).Msg("invalid signature")
		return nil, ErrInvalidToken
	}
	return token.Claims.(jwt.MapClaims), nil
}

func GetUserIDFromJWT(tokenString string) (int64, error) {
	claims, err := GetClaimsFromJWT(tokenString)
	if err != nil {
		return 0, err
	}
	userID, ok := claims[sub].(string)
	if !ok {
		log.Error().Fields(claims).Msg("error getting userID from JWT")
		return 0, ErrInvalidToken
	}
	i, err := strconv.ParseInt(userID, base, 64)
	if err != nil {
		log.Error().Str("userID", userID).Msg("error parsing userID to int64")
		return 0, ErrInvalidToken
	}
	return i, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrInvalidToken
	}
	return config.JWT, nil
}

func SignInToken(userID int64) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		sub: strconv.FormatInt(userID, base),
	}).SignedString(config.JWT)
	if err != nil {
		log.Error().Err(err).Int64("userID", userID).Msg("error creating sign in token")
		return "", err
	}
	log.Info().Int64("userID", userID).Msg("created sign in token")
	return token, nil
}
