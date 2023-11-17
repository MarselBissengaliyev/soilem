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
)

type UserPostgres struct {
	db *pgx.Conn
}

func NewUserPostgres(db *pgx.Conn) *UserPostgres {
	return &UserPostgres{db}
}

func (r *UserPostgres) Registration(user *model.User) (*model.User, error) {
	var foundUser *model.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existUser, _ := r.GetUserByUserName(foundUser.UserName)
	if existUser.UserName == foundUser.UserName {
		return nil, ErrUserAlreadyExists
	}

	if existUser.PhoneNumber == foundUser.PhoneNumber {
		return nil, ErrUserAlreadyExists
	}

	sql := fmt.Sprintf(
		`INSERT INTO %s (phone_number, password, full_name, user_name) VALUES ($1, $2, $3, $4)`,
		usersTable,
	)

	err := r.db.QueryRow(ctx, sql, user.PhoneNumber, user.Password, user.FullName, user.UserName).Scan(&foundUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert user")
	}

	return foundUser, nil
}

func (r *UserPostgres) GetUserByUserName(userName model.UserName) (*model.User, error) {
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

func (r *UserPostgres) GetUsers(searchTerm string, limit int) ([]*model.User, error) {
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
