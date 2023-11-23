package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type SMSCodePostgres struct {
	db *pgx.Conn
}

func NewSMSCodePostgres(db *pgx.Conn) *SMSCodePostgres {
	return &SMSCodePostgres{db}
}

func (r *SMSCodePostgres) SetCode(updateSMSCode model.SMSCode, userName model.UserName) (*model.SMSCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var smsCode *model.SMSCode

	sql := fmt.Sprintf(
		"UPDATE %s SET code=$1, expires_at=$2 WHERE username=$3",
		smsCodeTable,
	)
	if err := r.db.QueryRow(ctx, sql, updateSMSCode.Code, updateSMSCode.ExpiresAt).Scan(&smsCode); err != nil {
		return nil, errors.Wrap(err, "failed to update sms code")
	}

	return smsCode, nil
}
