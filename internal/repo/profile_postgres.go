package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type ProfilePostgres struct {
	db *pgx.Conn
}

func NewProfilePostgres(db *pgx.Conn) *ProfilePostgres {
	return &ProfilePostgres{db}
}

func (r *ProfilePostgres) Create(profile *model.Profile) (*model.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf(
		"INSERT INTO %s (f_name, l_name, description, avatar, date_of_birth, sex, author) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		profilesTable,
	)

	var createdProfile *model.Profile

	err := r.db.QueryRow(
		ctx, sql, profile.FName, profile.LName, profile.Description,
		profile.Avatar, profile.DateOfBirth, profile.Sex, profile.Author,
	).Scan(&createdProfile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create profile")
	}

	return createdProfile, nil
}

func (r *ProfilePostgres) GetByUserName(userName model.UserName) (*model.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf("SELECT * FROM %s WHERE user_name=$1", profilesTable)

	var profile *model.Profile

	if err := r.db.QueryRow(ctx, sql, userName).Scan(&profile); err != nil {
		return nil, errors.Wrap(err, "failed to query profile")
	}

	return profile, nil
}
