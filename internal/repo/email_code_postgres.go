package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type EmailCodePostgres struct {
	db *pgx.Conn
}

func NewEmailCodePostgres(db *pgx.Conn) *EmailCodePostgres {
	return &EmailCodePostgres{db}
}

func (r *EmailCodePostgres) SetCode(
	updateEmailCode model.EmailCode,
) (*model.EmailCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var emailCode *model.EmailCode

	sql := fmt.Sprintf("UPDATE %s SET code=$1, expires_at=$2 WHERE user_name=$3", emailCodeTable)
	
	err := r.db.QueryRow(ctx, sql, updateEmailCode.Code, updateEmailCode.ExpiresAt, updateEmailCode.UserName).Scan(&emailCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update email code")
	}

	return emailCode, nil
}
