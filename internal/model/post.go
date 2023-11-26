package model

import "time"

type Post struct {
	ID           uint      `json:"-"`
	Author       UserName  `json:"author"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title" validate:"required"`
	Description  string    `json:"description"`
	Content      string    `json:"content" validate:"required"`
	Location     string    `json:"location"`
	PrivacyLevel string    `json:"privacy_level" validate:"required"`
	ViewsCount   int       `json:"views_count"`
	HashTags     []string  `json:"hash_tags"`
	Images       []string  `json:"images"`
	Attachments  []string  `json:"attachments"`
	Links        []string  `json:"links"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *Post) Validate() error {
	return v.Struct(p)
}
