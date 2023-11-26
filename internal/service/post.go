package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
	"github.com/gosimple/slug"
)

type PostService struct {
	repo repo.Post
}

func NewPostService(repo repo.Post) *PostService {
	return &PostService{repo}
}

func (s *PostService) Create(post *model.Post, userName model.UserName) (*model.Post, *model.Fail) {
	if err := post.Validate(); err != nil {
		return nil, &model.Fail{
			Message:    "failed to validate profile: " + err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	createdAt := time.Now().UTC()
	post.CreatedAt = createdAt
	post.UpdatedAt = createdAt

	slugText := slug.Make(post.Title)
	existingPost, _ := s.repo.GetBySlug(slugText)

	counter := 2
	for existingPost != nil {
		// Добавляем счетчик к существующему слагу
		slugText = slug.Make(fmt.Sprintf("%s-%d", post.Title, counter))
		counter++

		existingPost, _ = s.repo.GetBySlug(slugText)
	}

	post.Slug = slugText
	post.Author = userName

	createdPost, err := s.repo.Create(post)
	if err != nil {
		return nil, &model.Fail{
			Message:    "failed to create post: %s" + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return createdPost, nil
}
