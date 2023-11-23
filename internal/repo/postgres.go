package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const (
	usersTable     = "users"
	smsCodeTable   = "sms_codes"
	emailCodeTable = "email_codes"
	profilesTable  = "profiles"
)

type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
}

func NewPostgresDB(cfg *Config, ctx context.Context) (*pgx.Conn, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
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
