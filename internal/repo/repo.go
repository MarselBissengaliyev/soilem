package repo

import (
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
)

type User interface {
	Registration(user *model.User) (*model.User, error)
	GetUserByUserName(userName model.UserName) (*model.User, error)
	GetUsers(searchTerm string, limit int) ([]*model.User, error)
}

type Repository struct {
	User
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		User: NewUserPostgres(db),
	}
}
