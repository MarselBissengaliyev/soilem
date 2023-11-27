package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type SessionPostgres struct {
	db *pgx.Conn
}

func NewAccessTokenPostgres(db *pgx.Conn) *SessionPostgres {
	return &SessionPostgres{db}
}

func (r *SessionPostgres) Create(accessToken *model.AccessToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf(
		"INSERT INTO %s (token, user_name, expires_at, created_at) VALUES ($1, $2, $3, $4)",
		accessTokensTable,
	)
	_, err := r.db.Exec(
		ctx, sql, accessToken.Token,
		accessToken.UserName, accessToken.ExpiresAt,
		accessToken.CreatedAt,
	)

	if err != nil {
		return errors.Wrap(err, "failed to insert session")
	}

	return nil
}

func (r *SessionPostgres) GetByToken(token string) (*model.Session, error) {
	var foundSession *model.Session

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf("SELECT * FROM %s WHERE token = $1", accessTokensTable)
	if err := r.db.QueryRow(ctx, sql, token).Scan(&foundSession); err != nil {
		return nil, errors.Wrap(err, "failed to get session by token")
	}

	return foundSession, nil
}

func (r *SessionPostgres) RemoveByToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf("DELETE FROM %s WHERE token = $1", accessTokensTable)

	_, err := r.db.Exec(ctx, sql, token)
	if err != nil {
		return errors.Wrap(err, "failed to delete session by token")
	}

	return nil
}
