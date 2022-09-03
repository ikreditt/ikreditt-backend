package models

import (
	"time"

	"github.com/fluffy-octo/ik-reddit-backend/graph/model"
	uuid "github.com/satori/go.uuid"
)

type Loan struct {
	Base
	UserID            uuid.UUID `gorm:"type:uuid"`
	Amount            int
	PaidAmount        int
	UserName          string
	ExpectedPayAmount int
	ExpectedPayDate   time.Time
	Status            model.LoanStatus
	Description       string
	CompletedDate     time.Time `gorm:"default:0"`
	PhoneNumber       string
	TransactionRef    string
	Payments          []Payments `gorm:"foreignKey:LoanID"`
}

type Payments struct {
	Base
	LoanID         uuid.UUID `gorm:"type:uuid"`
	PaidAmount     int
	TransactionRef string
}

func (l Loan) CreateToGraphData() *model.Loan {

	payments := func(b []Payments) []*model.LoanPayment {
		gqlPayments := make([]*model.LoanPayment, len(l.Payments))
		for idx, pay := range l.Payments {
			gqlPayments[idx] = &model.LoanPayment{
				PaymentID:  pay.ID.String(),
				LoanID:     pay.LoanID.String(),
				PaidAmount: pay.PaidAmount,
			}
		}
		return gqlPayments
	}(l.Payments)

	return &model.Loan{
		LoanID:            l.ID.String(),
		UserID:            l.UserID.String(),
		Amount:            l.Amount,
		PaidAmount:        l.PaidAmount,
		LoanDate:          l.CreatedAt,
		ExpectedPayDate:   l.ExpectedPayDate,
		Status:            l.Status,
		CompletedDate:     l.CompletedDate,
		PhoneNumber:       l.PhoneNumber,
		Payments:          payments,
		UserName:          l.UserName,
		ExpectedPayAmount: l.ExpectedPayAmount,
	}
}
