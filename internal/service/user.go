package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type UserService struct {
	repo repo.User
}

func NewUserService(repo repo.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Registration(user *model.User) (*model.User, *model.Fail) {
	user.IsRegistration = true

	if err := user.Validate(); err != nil {
		return nil, &model.Fail{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	hashedPassword, err := user.HashPassword()
	if err != nil {
		return nil, &model.Fail{
			Message:    "failed to hash password: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	user.Password = hashedPassword
	user.CreatedAt = time.Now().UTC()
	user.LastLoginAt = time.Now().UTC()
	user.GenerateConfirmationCode()

	foundUser, err := s.repo.Registration(user)
	if err != nil || foundUser != nil {
		if err == repo.ErrUserAlreadyExists {
			return nil, &model.Fail{
				Message:    err.Error(),
				StatusCode: http.StatusConflict,
			}
		}

		return nil, &model.Fail{
			Message:    "failed to register user: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return foundUser, nil
}

func (s *UserService) Login(user *model.User) (*model.User, *model.Fail) {
	if err := user.Validate(); err != nil {
		return nil, &model.Fail{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	user.LastLoginAt = time.Now().UTC()

	foundUser, err := s.repo.GetUserByUserName(user.UserName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.Fail{
				Message:    err.Error(),
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &model.Fail{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !foundUser.CheckPasswordHash(user.Password) {
		return nil, &model.Fail{
			Message:    "provided password does not match password hash",
			StatusCode: http.StatusBadRequest,
		}
	}

	return foundUser, nil
}

func (s *UserService) GetUserByUserName(userName model.UserName) (*model.User, *model.Fail) {
	if userName == "" {
		return nil, &model.Fail{
			Message:    "user_name field cannot is empty",
			StatusCode: http.StatusBadRequest,
		}
	}

	foundUser, err := s.repo.GetUserByUserName(userName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.Fail{
				Message:    "user not found: " + err.Error(),
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &model.Fail{
			Message:    "failed to get user: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return foundUser, nil
}

func (s *UserService) GetUsers(searchTerm string, limit string) ([]*model.User, *model.Fail) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, &model.Fail{
			Message:    "limit must be an integer: " + err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	if limitInt > 20 {
		return nil, &model.Fail{
			Message:    "limit must be greater than 20: " + err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	if limitInt < 1 {
		return nil, &model.Fail{
			Message:    "limit must be less than 20: " + err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	users, err := s.repo.GetUsers(searchTerm, limitInt)
	if err != nil {
		return nil, &model.Fail{
			Message:    "failed to get users: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return users, nil
}