package main

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
)

func AuthenticationMiddleware(ctx iris.Context) {
	token, err := ctx.Request().Cookie("token")
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail("No token found in cookie: "+err.Error()).
			Status(iris.StatusUnauthorized))
		return
	}

	signingKey := []byte(os.Getenv("JWT_SECRET"))
	tokenClaims, err := jwt.Parse(token.Value, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil || !tokenClaims.Valid {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail("Invalid token: "+err.Error()).
			Status(iris.StatusUnauthorized))
		return
	}

	claims := tokenClaims.Claims.(jwt.MapClaims)
	// for testing purposes
	fmt.Println(claims["sub"], claims["exp"])
	ctx.Values().Set("claims", claims)
	ctx.Next()
}
