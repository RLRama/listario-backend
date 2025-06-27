package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/RLRama/listario-backend/logger"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func SetupJWT() (*jwt.Signer, *jwt.Verifier, error) {
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return nil, nil, fmt.Errorf("JWT_SECRET_KEY environment variable not set")
	}
	jwtSecretKey := []byte(jwtSecret)

	expHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil || expHours <= 0 {
		expHours = 72
	}
	expiration := time.Hour * time.Duration(expHours)

	signer := jwt.NewSigner(jwt.HS256, jwtSecretKey, expiration)
	verifier := jwt.NewVerifier(jwt.HS256, jwtSecretKey)
	verifier.WithDefaultBlocklist()
	logger.Info().Msgf("JWT setup complete with expiration of %d hours", expHours)
	return signer, verifier, nil
}
