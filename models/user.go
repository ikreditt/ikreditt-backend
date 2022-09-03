package models

import (
	"github.com/fluffy-octo/ik-reddit-backend/graph/model"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Base
	Name            string
	Phone           string `gorm:"unique"`
	Email           string `gorm:"unique"`
	Pin             string
	HasPin          bool
	TotalLoan       int
	TotalPaid       int
	TimesLoaned     int
	Password        string
	LoanLimit       int
	CreditScore     int
	ProfilePhotoURL string
	FrontIDPhotoURL string
	HasCompletedKYC bool
	HasSubmittedKYC bool
	Loans           []Loan        `gorm:"foreignKey:UserID"`
	UserDetails     []UserDetails `gorm:"foreignKey:UserID"`
}

type UserDetails struct {
	Base
	UserID      uuid.UUID `gorm:"type:uuid"`
	Description string
	Answer      string
}

func (u User) CreateToGraphData() *model.User {
	details := func(d []UserDetails) []*model.UserDetails {
		gqlDetails := make([]*model.UserDetails, len(d))
		for idx, dtl := range d {
			gqlDetails[idx] = &model.UserDetails{
				Description: dtl.Description,
				Answer:      dtl.Answer,
			}
		}
		return gqlDetails
	}(u.UserDetails)

	loans := func(l []Loan) []*model.Loan {
		gqlLoans := make([]*model.Loan, len(l))
		for idx, loan := range l {
			gqlLoans[len(gqlLoans)-1-idx] = loan.CreateToGraphData()
		}
		return gqlLoans
	}(u.Loans)

	return &model.User{
		UserID:          u.ID.String(),
		Name:            u.Name,
		Phone:           u.Phone,
		LoanLimit:       u.LoanLimit,
		Email:           u.Email,
		TotalLoan:       u.TotalLoan,
		TotalPaid:       u.TotalPaid,
		TimesLoaned:     u.TimesLoaned,
		CreditScore:     u.CreditScore,
		CreatedAt:       u.CreatedAt,
		ProfilePhotoURL: u.ProfilePhotoURL,
		FrontIDPhotoURL: u.FrontIDPhotoURL,
		HasCompletedKyc: u.HasCompletedKYC,
		HasPin:          u.HasPin,
		HasSubmittedKyc: u.HasSubmittedKYC,
		Kycdetails:      details,
		Loans:           loans,
	}
}
