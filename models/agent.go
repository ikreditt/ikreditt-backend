package models

import "github.com/fluffy-octo/ik-reddit-backend/graph/model"

type Agent struct {
	Base
	Name       string
	Email      string `gorm:"index;unique;not null;"`
	Phone      string
	Address    string
	Password   string
	NationalID string
}

func (a Agent) CreateToGraphData() *model.Agent {
	return &model.Agent{
		ID:         a.ID.String(),
		Name:       a.Name,
		Email:      a.Email,
		Phone:      a.Phone,
		NationalID: a.NationalID,
		CreatedAt: a.CreatedAt,
		Address:    a.Address,
	}
}