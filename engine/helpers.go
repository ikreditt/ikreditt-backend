package engine

import (
	"errors"

	"github.com/fluffy-octo/ik-reddit-backend/graph/model"
	"github.com/fluffy-octo/ik-reddit-backend/models"
	"github.com/fluffy-octo/ik-reddit-backend/utils"
	"gorm.io/gorm/clause"
)

func FetchUserByID(userId string) (*models.User, error) {
	var user models.User
	err := utils.DB.Preload(clause.Associations).Preload(clause.Associations).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FetchUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := utils.DB.Preload(clause.Associations).Preload(clause.Associations).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FetchUserByAuthToken(jwt string) (*models.User, error) {
	userId, err := utils.ValidateJWTForAuthId(jwt)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	var user models.User
	err = utils.DB.Preload(clause.Associations).Preload(clause.Associations).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FetchUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := utils.DB.Preload(clause.Associations).Preload(clause.Associations).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FetchUsers(page int) ([]models.User, error) {
	var users []models.User
	err := utils.DB.Preload(clause.Associations).Preload(clause.Associations).Limit(50).Offset(50 * (page - 1)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func FetchLoansByStatus(userId string, page int, status model.LoanStatus) ([]models.Loan, error) {
	var loans []models.Loan
	err := utils.DB.Preload(clause.Associations).Where("user_id = ?", userId).Where("status = ?", status).Limit(50).Offset(50 * (page - 1)).Find(&loans).Error
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func FetchAgentByAuthToken(jwt string) (*models.Agent, error) {
	agentId, err := utils.ValidateJWTForAuthId(jwt)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	var agent models.Agent
	err = utils.DB.Where("id = ?", agentId).First(&agent).Error
	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func FetchAgentByPhone(phone string) (*models.Agent, error) {
	var agent models.Agent
	err := utils.DB.Where("phone = ?", phone).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func FetchAdminByToken(jwt string) (*models.Admin, error) {
	adminId, err := utils.ValidateJWTForAuthId(jwt)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	var admin models.Admin
	err = utils.DB.Where("id = ?", adminId).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func FetchAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	err := utils.DB.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func FetchAgentByID(agentId string) (*models.Agent, error) {
	var agent models.Agent
	err := utils.DB.Where("id = ?", agentId).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func FetchAgentByEmail(email string) (*models.Agent, error) {
	var agent models.Agent
	err := utils.DB.Where("email = ?", email).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func FetchAgents(page int) ([]models.Agent, error) {
	var agents []models.Agent
	err := utils.DB.Limit(50).Offset(50 * (page - 1)).Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func FetchLoans(page int, status model.LoanStatus) ([]models.Loan, error) {
	var loans []models.Loan
	err := utils.DB.Preload(clause.Associations).Limit(50).Offset(50*(page-1)).Where("status = ?", status).Order("created_at DESC").Find(&loans).Error
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func FetchLoansByUserID(userId string) ([]models.Loan, error) {
	var loans []models.Loan
	err := utils.DB.Preload(clause.Associations).Where("user_id = ?", userId).Find(&loans).Error
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func FetchLoanByID(loanId string) (*models.Loan, error) {
	var loan models.Loan
	err := utils.DB.Preload(clause.Associations).Where("id = ?", loanId).First(&loan).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

