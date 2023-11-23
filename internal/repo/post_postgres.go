package repo

import "github.com/jackc/pgx/v5"

type PostPostgres struct {
	db *pgx.Conn
}

func NewPostPostgres(db *pgx.Conn) *PostPostgres {
	return &PostPostgres{db}
}

func (p *PostPostgres) Create() {
	
}