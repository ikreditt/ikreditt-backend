package loans

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fluffy-octo/ik-reddit-backend/engine"
	"github.com/fluffy-octo/ik-reddit-backend/graph/model"
	"github.com/fluffy-octo/ik-reddit-backend/models"
	"github.com/fluffy-octo/ik-reddit-backend/utils"
)

var interest = 1.36
var limit = 50000

func ApplyFORLoan(loanApplication model.ApplyLoanInput) (*model.Loan, error) {
	user, err := engine.FetchUserByAuthToken(loanApplication.Token)
	if err != nil {
		return nil, err
	}

	if !user.HasSubmittedKYC {
		return nil, errors.New("please complete the kyc process")
	}

	if !user.HasCompletedKYC {
		return nil, errors.New("your verification is still in progress, Try again later")
	}

	if loanApplication.ExpectedPaymentDate.Before(time.Now().Add(time.Duration(time.Hour * 24))) {
		return nil, errors.New("payment duration must be at least 1 day from now")
	}

	if loanApplication.Amount < 500 {
		return nil, errors.New("loan amount must be at least 500")
	}

	if loanApplication.Amount > user.LoanLimit {
		return nil, fmt.Errorf("loan amount must be less than %v", user.LoanLimit)
	}

	userLoans, err := engine.FetchLoansByUserID(user.ID.String())
	if err != nil {
		return nil, err
	}

	for _, loan := range userLoans {
		switch loan.Status {
		case model.LoanStatusCurrent, model.LoanStatusCreated, model.LoanStatusOverdue:
			return nil, errors.New("user has an unpaid loan")
		}
	}

	log.Println(int( float64(loanApplication.Amount) * interest))

	newLoan := models.Loan{
		UserID:            user.ID,
		Amount:            loanApplication.Amount,
		PaidAmount:        0,
		ExpectedPayDate:   loanApplication.ExpectedPaymentDate,
		ExpectedPayAmount: int( float64(loanApplication.Amount) *  interest),
		Status:            model.LoanStatusCreated,
		UserName:          user.Name,
		Description:       loanApplication.Reason,
		PhoneNumber:       user.Phone,
	}

	err = utils.DB.Create(&newLoan).Error
	if err != nil {
		return nil, err
	}

	return newLoan.CreateToGraphData(), nil
}

func FetchUserLoans(input model.FetchUserLoansInput) ([]*model.Loan, error) {
	user, err := engine.FetchUserByAuthToken(input.Token)
	if err != nil {
		_, err = engine.FetchAgentByAuthToken(input.Token)
		if err != nil {
			_, err = engine.FetchAdminByToken(input.Token)
			if err != nil {
				return nil, err
			}
		}
		user, err = engine.FetchUserByID(input.UserID)
		if err != nil {
			return nil, errors.New("invalid user id")
		}
	}

	loans, err := engine.FetchLoansByStatus(user.ID.String(), input.Page, input.Status)
	if err != nil {
		return nil, err
	}

	var graphLoans []*model.Loan
	for _, loan := range loans {
		graphLoans = append(graphLoans, loan.CreateToGraphData())
	}
	return graphLoans, nil
}

func FetchLoans(input model.LoansInput) ([]*model.Loan, error) {
	page := 1
	if input.Page > 0 {
		page = input.Page
	}

	loans, err := engine.FetchLoans(page, input.Status)
	if err != nil {
		return nil, err
	}

	var graphLoans []*model.Loan
	for _, loan := range loans {
		graphLoans = append(graphLoans, loan.CreateToGraphData())
	}
	return graphLoans, nil
}

func ApproveNewLoan(input model.ApproveLoanInput) (bool, error) {
	_, err := engine.FetchAdminByToken(input.Token)
	if err != nil {
		return false, err
	}

	loan, err := engine.FetchLoanByID(input.LoanID)
	if err != nil {
		return false, err
	}

	user, err := engine.FetchUserByID(loan.UserID.String())
	if err != nil {
		return false, err
	}

	if loan.Status != model.LoanStatusCreated {
		return false, errors.New("invalid loan status")
	}

	timeDiff := loan.ExpectedPayDate.Sub(loan.CreatedAt)

	err = utils.DB.Model(&models.Loan{}).Where("id = ?", loan.ID).Updates(&models.Loan{
		ExpectedPayDate: time.Now().Add(timeDiff),
		Status:          model.LoanStatusCurrent,
		TransactionRef:  input.TransactionRef,
	}).Error

	if err != nil {
		return false, err
	}

	err = utils.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&models.User{TotalLoan: user.TotalLoan + loan.Amount, TimesLoaned: user.TimesLoaned + 1}).Error
	if err != nil {
		return false, err
	}

	//code to send money to user
	return true, nil
}

func RejectLoan(input model.RejectLoanInput) (bool, error) {
	_, err := engine.FetchAdminByToken(input.Token)
	if err != nil {
		return false, err
	}
	loan, err := engine.FetchLoanByID(input.LoanID)
	if err != nil {
		return false, err
	}

	if loan.Status != model.LoanStatusCreated {
		return false, errors.New("invalid loan status")
	}

	err = utils.DB.Model(&models.Loan{}).Where("id = ?", loan.ID).Updates(&models.Loan{
		Status: model.LoanStatusDeclined,
	}).Error

	if err != nil {
		return false, err
	}
	return true, nil
}

func PayLoan(input model.LoanPaymentInput) (bool, error) {
	//safaricom webhook example https://developer.safaricom.co.ke/docs/paybill-api/paybill-api-quick-start/
	user, err := engine.FetchUserByAuthToken(input.Token)
	if err != nil {
		return false, err
	}
	loan, err := engine.FetchLoanByID(input.LoanID)
	if err != nil {
		return false, err
	}

	if loan.PaidAmount >= loan.ExpectedPayAmount {
		return false, errors.New("loan already fully paid")
	}

	switch loan.Status {
	case model.LoanStatusCreated, model.LoanStatusDeclined, model.LoanStatusCompleted:
		return false, errors.New("invalid loan status")
	}

	if loan.PaidAmount+input.PayAmount > loan.ExpectedPayAmount {
		return false, errors.New("amount exceeds payment amount")
	}

	if loan.PaidAmount+input.PayAmount == loan.ExpectedPayAmount {
		err = utils.DB.Model(&models.Loan{}).Where("id = ?", loan.ID).Updates(&models.Loan{
			Status:        model.LoanStatusCompleted,
			CompletedDate: time.Now(),
		}).Error
		if err != nil {
			return false, err
		}

		fromLimit := limit - user.LoanLimit
		newLimit := ((loan.Amount / limit) * fromLimit) + user.LoanLimit

		err = utils.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&models.User{
			LoanLimit: newLimit,
		}).Error

		if err != nil {
			return false, err
		}
	}

	err = utils.DB.Model(&models.Loan{}).Where("id = ?", loan.ID).Updates(&models.Loan{
		PaidAmount: loan.PaidAmount + input.PayAmount,
	}).Error

	if err != nil {
		return false, err
	}

	newPayment := models.Payments{
		LoanID:         loan.ID,
		PaidAmount:     input.PayAmount,
		TransactionRef: input.TransactionRef,
	}

	err = utils.DB.Create(&newPayment).Error
	if err != nil {
		return false, err
	}

	err = utils.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&models.User{
		TotalPaid: user.TotalPaid + input.PayAmount,
	}).Error

	if err != nil {
		return false, err
	}

	return true, nil

}

func FetchLoanDetails(input model.FetchLoanInput) (*model.Loan, error) {
	_, err := engine.FetchUserByAuthToken(input.Token)
	if err != nil {
		_, err := engine.FetchAdminByToken(input.Token)
		if err != nil {
			_, err := engine.FetchAgentByAuthToken(input.Token)
			if err != nil {
				return nil, errors.New("invalid token")
			}
		}
	}
	loan, err := engine.FetchLoanByID(input.LoanID)
	if err != nil {
		return nil, err
	}
	return loan.CreateToGraphData(), nil
}
