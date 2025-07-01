package handler

import (
	"errors"

	"github.com/RLRama/listario-backend/logger"
	"github.com/RLRama/listario-backend/models"
	"github.com/RLRama/listario-backend/repository"
	"github.com/RLRama/listario-backend/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

type UserHandler struct {
	userService service.UserService
	verifier    *jwt.Verifier
}

func NewUserHandler(us service.UserService, verifier *jwt.Verifier) *UserHandler {
	return &UserHandler{
		userService: us,
		verifier:    verifier,
	}
}

func (h *UserHandler) Register(ctx iris.Context) {
	var req models.RegisterRequest
	if err := ctx.ReadJSON(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to read registration request")
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request format"})
		return
	}

	user, err := h.userService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			ctx.StatusCode(iris.StatusConflict)
			ctx.JSON(iris.Map{"error": err.Error()})
			return
		}
		logger.Error().Err(err).Msg("Failed to register user")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to register user"})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(user)
}

func (h *UserHandler) Login(ctx iris.Context) {
	var req models.LoginRequest
	if err := ctx.ReadJSON(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to read login request")
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "invalid request format"})
		return
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(iris.Map{"error": err.Error()})
			return
		}
		logger.Error().Err(err).Msg("User login failed")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "login failed"})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"token": token})
}

func (h *UserHandler) GetMyDetails(ctx iris.Context) {
	claims := jwt.Get(ctx).(*models.UserClaims)

	user, err := h.userService.GetUserDetails(claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": err.Error()})
			return
		}
		logger.Error().Err(err).Uint("userID", claims.UserID).Msg("Failed to get user details")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "could not retrieve user details"})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(user)
}

func (h *UserHandler) UpdateMyDetails(ctx iris.Context) {
	claims := jwt.Get(ctx).(*models.UserClaims)

	var req models.UpdateUserRequest
	if err := ctx.ReadJSON(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to read update request")
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "invalid request format"})
		return
	}

	user, err := h.userService.UpdateUserDetails(claims.UserID, req.Username, req.Email)
	if err != nil {
		logger.Error().Err(err).Uint("userID", claims.UserID).Msg("Failed to update user details")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "could not update user details"})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(user)
}

func (h *UserHandler) Logout(ctx iris.Context) {
	verifiedToken := jwt.GetVerifiedToken(ctx)
	if verifiedToken == nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "invalid token"})
		return
	}

	err := h.verifier.Blocklist.InvalidateToken(verifiedToken.Token, verifiedToken.StandardClaims)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to invalidate token")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "could not logout"})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "logout successful"})
}
