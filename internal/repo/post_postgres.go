package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type PostPostgres struct {
	db *pgx.Conn
}

func NewPostPostgres(db *pgx.Conn) *PostPostgres {
	return &PostPostgres{db}
}

func (r *PostPostgres) Create(post *model.Post) (*model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf(`
		INSERT INTO %s (
			author, slug, title, 
			description, content, location, 
			privacy_level, views_count, hash_tags, images, attachments, links, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, postsTable)

	var createdPost *model.Post

	if err := r.db.QueryRow(ctx, sql,
		post.Author, post.Slug, post.Title, post.Description,
		post.Content, post.Location, post.PrivacyLevel, post.ViewsCount,
		post.HashTags, post.Images, post.Attachments, post.Links,
		post.CreatedAt, post.UpdatedAt,
	).Scan(&createdPost); err != nil {
		return nil, errors.Wrap(err, "failed to create post")
	}

	return createdPost, nil
}

func (r *PostPostgres) GetBySlug(slug string) (*model.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf("SELECT * FROM %s WHERE slug=$1", profilesTable)

	var post *model.Post
	if err := r.db.QueryRow(ctx, sql, slug).Scan(&post); err != nil {
		return nil, errors.Wrap(err, "failed to get post by slug")
	}

	return post, nil
}