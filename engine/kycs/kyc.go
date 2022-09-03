package kycs

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/fluffy-octo/ik-reddit-backend/engine"
	"github.com/fluffy-octo/ik-reddit-backend/graph/model"
	"github.com/fluffy-octo/ik-reddit-backend/models"
	"github.com/fluffy-octo/ik-reddit-backend/utils"
)

var NewQA = []*model.KycQuestion{
	{QuestionText: "First Name", QuestionScore: 50},
	{QuestionText: "Second Name", QuestionScore: 50},
	{QuestionText: "Date of Birth", QuestionScore: 100},
	{QuestionText: "ID Number", QuestionScore: 100},
	{QuestionText: "Gender",
		Choices: []*model.KycChoice{
			{
				ChoiceName:  "Male",
				ChoiceScore: 100,
			}, {
				ChoiceName:  "Female",
				ChoiceScore: 100,
			},
		},
	},
	{QuestionText: "Alternate Contact No", QuestionScore: 100},
	{QuestionText: "Home Address", QuestionScore: 100},
	{QuestionText: "Are you a breadwinner?/Breadwinner",
		Choices: []*model.KycChoice{
			{
				ChoiceName:  "Yes",
				ChoiceScore: 100,
			}, {
				ChoiceName:  "No",
				ChoiceScore: 0,
			},
		},
	},
	{QuestionText: "HouseHold Information",
		Choices: []*model.KycChoice{
			{
				ChoiceName:  "Rental",
				ChoiceScore: 100,
			}, {
				ChoiceName:  "Family Owned",
				ChoiceScore: 50,
			}, {
				ChoiceName:  "Owned",
				ChoiceScore: 150,
			},
		},
	},
	{QuestionText: "Nature of work",
		Choices: []*model.KycChoice{
			{
				ChoiceName:  "Fulltime",
				ChoiceScore: 200,
			}, {
				ChoiceName:  "Part-time",
				ChoiceScore: 100,
			}, {
				ChoiceName:  "Business owner",
				ChoiceScore: 150,
			},
			{
				ChoiceName:  "Self-employed",
				ChoiceScore: 100,
			}, {
				ChoiceName:  "Student",
				ChoiceScore: 50,
			},
			{
				ChoiceName:  "Unemployed",
				ChoiceScore: 0,
			},
		},
	},
	{QuestionText: "What industry is your work?/Job Industry"},
	{QuestionText: "How much do you earn?/Salary Range",
		Choices: []*model.KycChoice{
			{
				ChoiceName:  "0 - 5,000",
				ChoiceScore: 50,
			}, {
				ChoiceName:  "5,001 - 10,000",
				ChoiceScore: 50,
			}, {
				ChoiceName:  "10,001 - 20,000",
				ChoiceScore: 100,
			},
			{
				ChoiceName:  "20,001 - 35,000",
				ChoiceScore: 100,
			}, {
				ChoiceName:  "Above 35,000",
				ChoiceScore: 150,
			},
		},
	},

	{QuestionText: "How frequently do you get paid?/Income Frequency",
		Choices: []*model.KycChoice{
			{
				ChoiceName:  "Everyday",
				ChoiceScore: 50,
			}, {
				ChoiceName:  "Every week",
				ChoiceScore: 100,
			}, {
				ChoiceName:  "Twice each week",
				ChoiceScore: 150,
			},
			{
				ChoiceName:  "Monthly",
				ChoiceScore: 200,
			},
		},
	},

	{QuestionText: "Select a date which you expect your next pay"},

	{QuestionText: "Level of Education",
		Choices: []*model.KycChoice{
			{
				ChoiceName:  "Primary",
				ChoiceScore: 0,
			}, {
				ChoiceName:  "High school",
				ChoiceScore: 50,
			}, {
				ChoiceName:  "University",
				ChoiceScore: 100,
			},
			{
				ChoiceName:  "Masters/Phd",
				ChoiceScore: 100,
			},
		},
	},
	{QuestionText: "Referee 1 Name", QuestionScore: 10},
	{QuestionText: "Referee 1 Phone", QuestionScore: 20},
	{QuestionText: "Referee 2 Name", QuestionScore: 10},
	{QuestionText: "Referee 2 Phone", QuestionScore: 20},
}

func FetchKycs(input model.FetchKycInput) ([]*model.KycQuestion, error) {
	_, err := engine.FetchUserByAuthToken(input.Token)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return NewQA, nil
}

func SaveUserKyc(input model.KycDetails) (bool, error) {
	user, err := engine.FetchUserByAuthToken(input.Token)
	if err != nil {
		return false, err
	}

	if user.HasSubmittedKYC {
		return false, errors.New("user has already submitted KYC")
	}

	if user.HasCompletedKYC {
		return false, errors.New("user has already completed KYC")
	}

	found := 0
	for _, d := range input.Details {
		for _, realq := range NewQA {
			if d.Description == realq.QuestionText && d.Answer != "" {
				found++
			}
		}
	}

	if found < len(NewQA) {
		return false, errors.New("missing required details")
	}

DetailsLoop:
	for _, d := range input.Details {
		for _, saved := range user.UserDetails {
			if saved.Description == d.Description && saved.Answer != "" {
				continue DetailsLoop
			}
		}

		var newDetail = models.UserDetails{
			UserID:      user.ID,
			Description: d.Description,
			Answer:      d.Answer,
		}
		err = utils.DB.Create(&newDetail).Error
		if err != nil {
			return false, err
		}
	}

	// db.Model(&models.Merchants{}).Where("email = ?", email).Updates(&models.Merchants{
	err = utils.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&models.User{HasSubmittedKYC: true}).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func ApproveUserKYC(input model.VerifyKycInput) (bool, error) {
	_, err := engine.FetchAgentByAuthToken(input.Token)
	if err != nil {
		_, err := engine.FetchAdminByToken(input.Token)
		if err != nil {
			return false, err
		}
	}

	user, err := engine.FetchUserByID(input.UserID)
	if err != nil {
		return false, errors.New("user not found")
	}

	if !user.HasSubmittedKYC {
		return false, errors.New("user has not submitted KYC")
	}

	if user.HasCompletedKYC {
		return false, errors.New("user has already completed KYC")
	}

	score, err := ComputeKycScore(input.UserID)
	if err != nil {
		return false, err
	}

	err = utils.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&models.User{CreditScore: score, LoanLimit: score * 10, HasCompletedKYC: true}).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func ComputeKycScore(userId string) (int, error) {
	user, err := engine.FetchUserByID(userId)
	if err != nil {
		return 0, err
	}

	var score int
	for _, d := range user.UserDetails {
		for _, q := range NewQA {
			if d.Description == q.QuestionText {
				if q.Choices == nil {
					if d.Answer != "" {
						score += q.QuestionScore
					}
				} else {
					for _, c := range q.Choices {
						if d.Answer == c.ChoiceName {
							score += c.ChoiceScore
						}
					}
				}
			}
		}
	}

	return score, nil
}

type IdVerificationResponse struct {
	Reason   string `json:"reason"`
	UserInfo struct {
		IDNumber     int `json:"ID_NUMBER"`
		SerialNumber int `json:"SERIAL_NUMBER"`
		DateOfIssue  struct {
			Date           int   `json:"date"`
			Day            int   `json:"day"`
			Hours          int   `json:"hours"`
			Minutes        int   `json:"minutes"`
			Month          int   `json:"month"`
			Nanos          int   `json:"nanos"`
			Seconds        int   `json:"seconds"`
			Time           int64 `json:"time"`
			TimezoneOffset int   `json:"timezoneOffset"`
			Year           int   `json:"year"`
		} `json:"DATE_OF_ISSUE"`
		FullNames       string `json:"FULL_NAMES"`
		Gender          string `json:"GENDER"`
		DateOfBirth     int    `json:"DATE_OF_BIRTH"`
		DistrictOfBirth string `json:"DISTRICT_OF_BIRTH"`
		FatherNames     string `json:"FATHER_NAMES"`
		MotherNames     string `json:"MOTHER_NAMES"`
		Address         string `json:"ADDRESS"`
	} `json:"userInfo"`
	Result int `json:"result"`
}

func VerifyIDBySerialNumber(serialNo string) (bool, error) {
	var timsurl = "https://tims.ntsa.go.ke/rbac/user/getIsIDRegistered.htm?"

	var res *http.Response

	endpoint, err := url.Parse(timsurl)
	if err != nil {
		return false, err
	}
	params := endpoint.Query()

	params.Set("idNo", serialNo)

	endpoint.RawQuery = params.Encode()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err = client.Get(endpoint.String())
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var response IdVerificationResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	if response.Result == 1 {
		return false, errors.New(response.Reason)
	}

	return true, nil
}
