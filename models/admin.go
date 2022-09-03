package models

type Admin struct {
	Base
	Email    string `gorm:"uniqueIndex"`
	Password string
}
