package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/RLRama/listario-backend/logger"
	"github.com/kataras/iris/v12/middleware/jwt"
)

const (
	accessTokenMaxAge  = 15 * time.Minute
	refreshTokenMaxAge = 7 * 24 * time.Hour
)

func SetupJWT() (*jwt.Signer, *jwt.Verifier, time.Duration, error) {
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return nil, nil, 0, fmt.Errorf("JWT_SECRET_KEY environment variable not set")
	}
	jwtSecretKey := []byte(jwtSecret)

	signer := jwt.NewSigner(jwt.HS256, jwtSecretKey, accessTokenMaxAge)

	verifier := jwt.NewVerifier(jwt.HS256, jwtSecretKey)
	verifier.WithDefaultBlocklist()

	logger.Info().Msgf("JWT setup complete. Access token lifespan: %s", accessTokenMaxAge)
	return signer, verifier, refreshTokenMaxAge, nil
}
