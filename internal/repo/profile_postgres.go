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

func (r *ProfilePostgres) Create(profile *model.Profile) (*model.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf(
		"INSERT INTO %s (f_name, l_name, description, avatar, date_of_birth, sex) VALUES ($1, $2, $3, $4, $5, $6)",
		profilesTable,
	)

	var createdProfile *model.Profile

	err := r.db.QueryRow(
		ctx, sql, profile.FName, profile.LName, profile.Description,
		profile.Avatar, profile.DateOfBirth, profile.Sex,
	).Scan(&createdProfile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create profile")
	}

	return createdProfile, nil
}
