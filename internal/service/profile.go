package service

import (
	"fmt"
	"net/http"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
)

type ProfileService struct {
	repo repo.Profile
}

func NewProfileService(repo repo.Profile) *ProfileService {
	return &ProfileService{repo}
}

func (s *ProfileService) Create(profile *model.Profile) (*model.Profile, *model.Fail) {
	if err := profile.Validate(); err != nil {
		return nil, &model.Fail{
			Message:    "failed to validate profile: " + err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	existingProfile, _ := s.repo.GetByUserName(profile.Author)
	if existingProfile != nil {
		return nil, &model.Fail{
			Message:    fmt.Sprintf("profile with user_name %s already exists", profile.Author),
			StatusCode: http.StatusConflict,
		}
	}

	createdProfile, err := s.repo.Create(profile)
	if err != nil {
		return nil, &model.Fail{
			Message:    "failed to create profile: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return createdProfile, nil
}
