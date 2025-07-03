package service

import (
	"errors"

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
	Login(email, password string) (string, error)
	GetUserDetails(userID uint) (*models.User, error)
	UpdateUserDetails(userID uint, username, email string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	signer   *jwt.Signer
}

func NewUserService(repo repository.UserRepository, signer *jwt.Signer) UserService {
	return &userService{
		userRepo: repo,
		signer:   signer,
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

func (s *userService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", ErrInvalidCredentials
	}

	claims := models.UserClaims{UserID: user.ID}

	token, err := s.signer.Sign(claims)
	if err != nil {
		return "", err
	}

	return string(token), nil
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
