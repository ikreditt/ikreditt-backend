package admins

import (
	"errors"
	"log"

	"github.com/fluffy-octo/ik-reddit-backend/engine"
	"github.com/fluffy-octo/ik-reddit-backend/graph/model"
	"github.com/fluffy-octo/ik-reddit-backend/models"
	"github.com/fluffy-octo/ik-reddit-backend/utils"
)

func LoginAdmin(loginValues model.LoginInput) (*model.AgentAuthPayload, error) {
	admin, err := validateAdminLogin(loginValues.Email, loginValues.Password)
	if err != nil {
		return nil, err
	}

	authenticationToken, err := utils.GenerateJWTForAuthId(&admin.ID)
	if err != nil {
		return nil, errors.New("couldn't generate token")
	}

	AuthPayload := model.AgentAuthPayload{
		Token: authenticationToken,
	}

	return &AuthPayload, nil
}

func validateAdminLogin(email, password string) (*models.Admin, error) {
	admin, err := engine.FetchAdminByEmail(email)

	if err != nil {
		return nil, errors.New("invalid email or/and password")
	}

	if !utils.CompareHashedString(admin.Password, password) {
		return nil, errors.New("invalid email or/and password")
	}
	return admin, nil
}

func AddAdmin() {
	var admin models.Admin

	err := utils.DB.First(&admin).Error
	if err == nil {
		log.Println("Skipping: admin already exists")
		return
	}

	admin.Email = "admin@test.com"
	
	password, err := utils.HashString("1krediTt")
	if err != nil {
		log.Fatal("internal encryption error")
	}

	admin.Password = password
	err = utils.DB.Save(&admin).Error
	if err != nil {
		return
	}
}
