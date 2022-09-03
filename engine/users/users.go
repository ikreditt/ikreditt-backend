package users

import (
	"errors"
	"fmt"
	"time"

	"github.com/fluffy-octo/ik-reddit-backend/engine"
	"github.com/fluffy-octo/ik-reddit-backend/graph/model"
	"github.com/fluffy-octo/ik-reddit-backend/models"
	"github.com/fluffy-octo/ik-reddit-backend/utils"
)

const firstLoanLimit = 500
const rollOver = 1.02

func RegisterUser(user model.UserRegisterInput) (*model.RegistertionPayload, error) {
	if !(user.Password == user.ConfirmPassword) {
		return nil, errors.New("passwords don't match")
	}

	_, err := engine.FetchUserByPhone(user.Phone)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	_, err = engine.FetchUserByEmail(user.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	password, err := utils.HashString(user.Password)
	if err != nil {
		return nil, errors.New("internal encryption error")
	}
	newUser := models.User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  password,
		Phone:     user.Phone,
		LoanLimit: firstLoanLimit,
	}

	err = utils.DB.Create(&newUser).Error
	if err != nil {
		return nil, errors.New("failed to Save Data")
	}

	authenticationToken, err := utils.GenerateJWTForAuthId(&newUser.ID)
	if err != nil {
		return nil, errors.New("couldn't generate token")
	}

	//TODO:send verification email
	resp := model.RegistertionPayload{
		Token: authenticationToken,
	}

	return &resp, nil
}

func SetUserPin(PinInput model.SetPinInput) (bool, error) {
	user, err := engine.FetchUserByAuthToken(PinInput.Token)
	if err != nil {
		return false, err
	}
	if user.HasPin {
		return false, errors.New("user already has pin")
	}
	pin, err := utils.HashString(fmt.Sprint(PinInput.Pin))
	if err != nil {
		return false, errors.New("internal encryption error")
	}

	err = utils.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&models.User{
		Pin:    pin,
		HasPin: true,
	}).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func LoginUserByPhone(loginValues model.UserLoginInput) (*model.RegistertionPayload, error) {
	user, err := validateUserLogin(loginValues.Phone, loginValues.Password)
	if err != nil {
		return nil, err
	}

	authenticationToken, err := utils.GenerateJWTForAuthId(&user.ID)
	if err != nil {
		return nil, errors.New("couldn't generate token")
	}

	resp := model.RegistertionPayload{
		Token: authenticationToken,
	}

	return &resp, nil
}

func validateUserLogin(phone, password string) (*models.User, error) {
	user, err := engine.FetchUserByPhone(phone)
	if err != nil {
		return nil, errors.New("invalid phone or/and password")
	}

	if !utils.CompareHashedString(user.Password, password) {
		return nil, errors.New("invalid phone or/and password")
	}
	return user, nil
}

func LoginUserByPin(loginValues model.PinLoginInput) (*model.RegistertionPayload, error) {
	user, err := validatePinLogin(fmt.Sprint(loginValues.Pin), loginValues.Token)
	if err != nil {
		return nil, err
	}
	authenticationToken, err := utils.GenerateJWTForAuthId(&user.ID)
	if err != nil {
		return nil, errors.New("couldn't generate token")
	}
	resp := model.RegistertionPayload{
		Token: authenticationToken,
	}
	return &resp, nil
}

func validatePinLogin(pin, token string) (*models.User, error) {
	user, err := engine.FetchUserByAuthToken(token)
	if err != nil {
		return nil, err
	}
	if !user.HasPin {
		return nil, errors.New("user doesn't have pin, Login by email instead")
	}
	if !utils.CompareHashedString(user.Pin, pin) {
		return nil, errors.New("wrong pin")
	}
	return user, nil
}

func FetchAllCustomers(input *model.PaginationInput) ([]*model.User, error) {
	page := 1
	if input != nil && input.Page > 0 {
		page = input.Page
	}

	customers, err := engine.FetchUsers(page)
	if err != nil {
		return nil, err
	}

	var graphCustomers []*model.User
	for _, customer := range customers {
		graphCustomers = append(graphCustomers, customer.CreateToGraphData())
	}
	return graphCustomers, nil
}

func FetchCustomerDetails(input model.FetchUserInput) (*model.User, error) {
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

	for _, loan := range user.Loans {
		if loan.Status == model.LoanStatusCurrent {
			if loan.ExpectedPayDate.Before(time.Now()) {
				dayspast := time.Since(loan.ExpectedPayDate).Hours() / 24
				err = utils.DB.Model(&models.Loan{}).Where("id = ?", loan.ID).Updates(&models.Loan{
					Status:            model.LoanStatusOverdue,
					ExpectedPayAmount: loan.ExpectedPayAmount * int(dayspast*rollOver),
				}).Error
				if err != nil {
					return nil, errors.New("failed to update loan status")
				}
			}
		}
	}

	return user.CreateToGraphData(), nil
}

func ChangePassword(input model.ChangePasswordInput) (bool, error) {
	user, err := engine.FetchUserByAuthToken(input.Token)
	if err != nil {
		return false, err
	}
	if !(input.NewPassword == input.ConfirmPassword) {
		return false, errors.New("passwords don't match")
	}
	if !utils.CompareHashedString(user.Password, input.OldPassword) {
		return false, errors.New("wrong old password")
	}
	if input.NewPassword == input.OldPassword {
		return false, errors.New("new password is the same as old password, Enter a different password instead")
	}
	newPassword, err := utils.HashString(input.NewPassword)
	if err != nil {
		return false, errors.New("internal encryption error")
	}

	err = utils.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&models.User{
		Password: newPassword,
	}).Error

	if err != nil {
		return false, err
	}
	return true, nil
}

func ChangePin(input model.ChangePinInput) (bool, error) {
	user, err := engine.FetchUserByAuthToken(input.Token)
	if err != nil {
		return false, err
	}
	if !user.HasPin {
		return false, errors.New("user doesn't have pin, Login by email instead")
	}
	if utils.CompareHashedString(user.Pin, fmt.Sprint(input.NewPin)) {
		return false, errors.New("you cannot set your old pin as your new pin, Enter a new pin instead")
	}
	if !utils.CompareHashedString(user.Password, input.Password) {
		return false, errors.New("wrong password")
	}
	newPin, err := utils.HashString(fmt.Sprint(input.NewPin))
	if err != nil {
		return false, errors.New("internal encryption error")
	}

	err = utils.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&models.User{
		Pin: newPin,
	}).Error

	if err != nil {
		return false, err
	}
	return true, nil
}

func SearchUserPhone(phone string) (string, error) {
	user, err := engine.FetchUserByPhone(phone)
	if err != nil {
		return "", errors.New("user not found")
	}
	return user.ID.String(), nil
}
