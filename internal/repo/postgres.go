package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const (
	usersTable        = "users"
	smsCodeTable      = "sms_codes"
	emailCodeTable    = "email_codes"
	profilesTable     = "profiles"
	postsTable        = "posts"
	accessTokensTable = "access_tokens"
)

type PostgresConfig struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
	SSLmode  string
}

func NewPostgresDB(cfg *PostgresConfig) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	db, err := pgx.Connect(
		ctx,
		connString,
	)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	return db, nil
}
