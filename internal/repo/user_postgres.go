package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserPostgres struct {
	db *pgx.Conn
}

func NewUserPostgres(db *pgx.Conn) *UserPostgres {
	return &UserPostgres{db}
}

func (r *UserPostgres) Create(user *model.User) (*model.User, error) {
	var createdUser *model.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf(
		`INSERT INTO %s (phone_number, password, user_name) VALUES ($1, $2, $4)`,
		usersTable,
	)

	tx, err := r.db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}

	err = tx.QueryRow(
		ctx, sql,
		user.PhoneNumber, user.Password, user.UserName,
	).Scan(&createdUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert user")
	}

	return createdUser, nil
}

func (r *UserPostgres) GetByUserName(userName model.UserName) (*model.User, error) {
	var foundUser *model.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf("SELECT * FROM %s WHERE user_name=$1", usersTable)

	err := r.db.QueryRow(ctx, sql, userName).Scan(&foundUser)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by user_name")
	}

	return foundUser, nil
}

func (r *UserPostgres) GetAll(searchTerm string, limit int) ([]*model.User, error) {
	var users []*model.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf("SELECT * FROM %s WHERE MATCH (full_name,user_name) AGAINST ('$1')", usersTable)
	rows, err := r.db.Query(ctx, sql, searchTerm, limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}

	defer rows.Close()

	if err := rows.Scan(&users); err != nil {
		return nil, errors.Wrap(err, "failed to scan rows in users")
	}

	return users, nil
}

func (r *UserPostgres) SetPhoneVerifiedValue(status bool, userName model.UserName) (bool, error) {
	var isPhoneVerified bool

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf(
		"UPDATE %s SET is_phone_verified=$1 WHERE username=$2 RETURNING is_phone_verified",
		usersTable,
	)
	if err := r.db.QueryRow(ctx, sql, status, userName).Scan(&isPhoneVerified); err != nil {
		return false, errors.Wrap(err, "failed to update is_phone_verified value")
	}

	return isPhoneVerified, nil
}

func (r *UserPostgres) SetEmailVerifiedValue(status bool, userName model.UserName) (bool, error) {
	var isEmailVerified bool

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sql := fmt.Sprintf(
		"UPDATE %s SET is_email_verified=$1 WHERE username=$2 RETURNING is_email_verified",
		usersTable,
	)
	if err := r.db.QueryRow(ctx, sql, status, userName).Scan(&isEmailVerified); err != nil {
		return false, errors.Wrap(err, "failed to update is_email_verified value")
	}

	return isEmailVerified, nil
}
