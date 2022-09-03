package agents

import (
	"errors"

	"github.com/fluffy-octo/ik-reddit-backend/engine"
	"github.com/fluffy-octo/ik-reddit-backend/engine/admins"
	"github.com/fluffy-octo/ik-reddit-backend/graph/model"
	"github.com/fluffy-octo/ik-reddit-backend/models"
	"github.com/fluffy-octo/ik-reddit-backend/utils"
)

func RegisterAgent(input model.RegisterAgentInput) (bool, error) {
	_, err := engine.FetchAdminByToken(input.Token)
	if err != nil {
		return false, err
	}

	_, err = engine.FetchAgentByEmail(input.Email)
	if err == nil {
		return false, errors.New("agent already exists")
	}
	_, err = engine.FetchAgentByPhone(input.Phone)
	if err == nil {
		return false, errors.New("agent already exists")
	}

	password, err := utils.HashString(input.Password)
	if err != nil {
		return false, errors.New("internal encryption error")
	}
	newAgent := models.Agent{
		Name:       input.Name,
		Email:      input.Email,
		Address:    input.Address,
		Password:   password,
		Phone:      input.Phone,
		NationalID: input.NationalID,
	}

	err = utils.DB.Create(&newAgent).Error
	if err != nil {
		return false, errors.New("failed to Save Data")
	}

	return true, nil
}

func LoginAgent(loginValues model.LoginInput) (*model.AgentAuthPayload, error) {
	agent, err := validateEmailLogin(loginValues.Email, loginValues.Password)
	if err != nil {
		return admins.LoginAdmin(loginValues)
	}

	authenticationToken, err := utils.GenerateJWTForAuthId(&agent.ID)
	if err != nil {
		return nil, errors.New("couldn't generate token")
	}

	authPayload := model.AgentAuthPayload{
		Token: authenticationToken,
	}
	return &authPayload, nil
}

func validateEmailLogin(email, password string) (*models.Agent, error) {
	agent, err := engine.FetchAgentByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or/and password")
	}

	if !utils.CompareHashedString(agent.Password, password) {
		return nil, errors.New("invalid email or/and password")
	}
	return agent, nil
}

func FetchAllAgents(input *model.PaginationInput) ([]*model.Agent, error) {
	page := 1
	if input != nil && input.Page >= 0 {
		page = input.Page
	}

	agents, err := engine.FetchAgents(page)
	if err != nil {
		return nil, errors.New("failed to fetch agents")
	}

	var agentList []*model.Agent
	for _, agent := range agents {
		agentList = append(agentList, agent.CreateToGraphData())
	}
	return agentList, nil
}

func FetchAgentProfile(input model.FetchAgentInput) (*model.Agent, error) {
	agent, err := engine.FetchAgentByAuthToken(input.Token)
	
	if err != nil {
		_, err = engine.FetchAdminByToken(input.Token)
		if err != nil {
			return nil, err
		}
		agent, err = engine.FetchAgentByID(input.ID)
		if err != nil {
			return nil, errors.New("agent not found")
		}
	}

	return agent.CreateToGraphData(), nil
}