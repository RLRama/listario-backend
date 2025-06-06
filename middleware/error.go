package middleware

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/RLRama/listario-backend/models"
)

func newErrorMiddleware() iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			if r := recover(); r != nil {
				if apiErr, ok := r.(*models.APIError); ok {
			
				}
			}
		}
	}
}