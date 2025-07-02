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

// Register
// @Summary      Register a new user
// @Description  Creates a new user account with a username, email, and password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload   body      models.RegisterRequest  true  "User Registration Payload"
// @Success      201  {object}  models.UserResponse
// @Failure      400  {object}  object{error=string} "Invalid request format"
// @Failure      409  {object}  object{error=string} "User with this email already exists"
// @Failure      500  {object}  object{error=string} "Failed to register user"
// @Router       /auth/register [post]
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

	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(response)
}

// Login
// @Summary      Log in a user
// @Description  Logs in a user with an email and password, returning a JWT.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload   body      models.LoginRequest     true  "User Login Payload"
// @Success      200  {object}  object{token=string}
// @Failure      400  {object}  object{error=string} "Invalid request format"
// @Failure      401  {object}  object{error=string} "Invalid credentials"
// @Failure      500  {object}  object{error=string} "Login failed"
// @Router       /auth/login [post]
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

// GetMyDetails
// @Summary      Get current user details
// @Description  Retrieves the details for the currently authenticated user.
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.UserResponse
// @Failure      401  {object}  object{error=string} "Unauthorized"
// @Failure      404  {object}  object{error=string} "User not found"
// @Failure      500  {object}  object{error=string} "Could not retrieve user details"
// @Router       /users/me [get]
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

	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(response)
}

// UpdateMyDetails
// @Summary      Update current user details
// @Description  Updates the username and/or email for the currently authenticated user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload   body      models.UpdateUserRequest  true  "User Update Payload"
// @Success      200  {object}  models.UserResponse
// @Failure      400  {object}  object{error=string} "Invalid request format"
// @Failure      401  {object}  object{error=string} "Unauthorized"
// @Failure      409  {object}  object{error=string} "User with this email already exists"
// @Failure      500  {object}  object{error=string} "Could not update user details"
// @Router       /users/me [put]
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
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			ctx.StatusCode(iris.StatusConflict)
			ctx.JSON(iris.Map{"error": err.Error()})
			return
		}

		logger.Error().Err(err).Uint("userID", claims.UserID).Msg("Failed to update user details")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "could not update user details"})
		return
	}

	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(response)
}

// Logout
// @Summary      Log out the current user
// @Description  Invalidates the current user's JWT, effectively logging them out.
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{message=string}
// @Failure      401  {object}  object{error=string} "Invalid token"
// @Failure      500  {object}  object{error=string} "Could not logout"
// @Router       /users/logout [get]
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
