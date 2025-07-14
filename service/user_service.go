package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/RLRama/listario-backend/models"
	"github.com/RLRama/listario-backend/repository"
	"github.com/RLRama/listario-backend/utils"
	"github.com/kataras/iris/v12/middleware/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (jwt.TokenPair, error)
	RefreshToken(userID uint) (jwt.TokenPair, error)
	GetUserDetails(userID uint) (*models.User, error)
	UpdateUserDetails(userID uint, username, email string) (*models.User, error)
}

type userService struct {
	userRepo           repository.UserRepository
	signer             *jwt.Signer
	refreshTokenMaxAge time.Duration
}

func NewUserService(repo repository.UserRepository, signer *jwt.Signer, refreshTokenMaxAge time.Duration) UserService {
	return &userService{
		userRepo:           repo,
		signer:             signer,
		refreshTokenMaxAge: refreshTokenMaxAge,
	}
}

func (s *userService) Register(username, email, password string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(email, password string) (jwt.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return jwt.TokenPair{}, ErrInvalidCredentials
		}
		return jwt.TokenPair{}, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return jwt.TokenPair{}, ErrInvalidCredentials
	}

	return s.generateTokenPair(user.ID)
}

func (s *userService) RefreshToken(userID uint) (jwt.TokenPair, error) {
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return jwt.TokenPair{}, err
	}
	return s.generateTokenPair(userID)
}

func (s *userService) generateTokenPair(userID uint) (jwt.TokenPair, error) {
	accessClaims := models.UserClaims{UserID: userID}
	refreshClaims := jwt.Claims{Subject: strconv.FormatUint(uint64(userID), 10)}

	return s.signer.NewTokenPair(accessClaims, refreshClaims, s.refreshTokenMaxAge)
}

func (s *userService) GetUserDetails(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) UpdateUserDetails(userID uint, username, email string) (*models.User, error) {
	if email != "" {
		existingUser, err := s.userRepo.FindByEmail(email)
		if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, repository.ErrUserAlreadyExists
		}
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
