package model

import "time"

type Profile struct {
	ID          uint      `json:"-"`
	Author      UserName  `json:"author"`
	FName       string    `json:"first_name" validate:"required"`
	LName       string    `json:"last_name" validate:"required"`
	Description string    `json:"description"`
	Avatar      string    `json:"avatar"`
	DateOfBirth time.Time `json:"birth_of_date" validate:"required"`
	Sex         string    `json:"sex" validate:"required, oneof=male female"`
}

func (p *Profile) Validate() error {
	return v.Struct(p)
}
