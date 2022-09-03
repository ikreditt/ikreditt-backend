package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/fluffy-octo/ik-reddit-backend/engine/agents"
	"github.com/fluffy-octo/ik-reddit-backend/engine/kycs"
	"github.com/fluffy-octo/ik-reddit-backend/engine/loans"
	"github.com/fluffy-octo/ik-reddit-backend/engine/users"
	"github.com/fluffy-octo/ik-reddit-backend/graph/generated"
	"github.com/fluffy-octo/ik-reddit-backend/graph/model"
)

func (r *mutationResolver) ApproveLoan(ctx context.Context, input model.ApproveLoanInput) (bool, error) {
	return loans.ApproveNewLoan(input)
}

func (r *mutationResolver) RejectLoan(ctx context.Context, input model.RejectLoanInput) (bool, error) {
	return loans.RejectLoan(input)
}

func (r *mutationResolver) RegisterAgent(ctx context.Context, input model.RegisterAgentInput) (bool, error) {
	return agents.RegisterAgent(input)
}

func (r *mutationResolver) Register(ctx context.Context, input model.UserRegisterInput) (*model.RegistertionPayload, error) {
	return users.RegisterUser(input)
}

func (r *mutationResolver) VerifyIDSerialNumber(ctx context.Context, input string) (bool, error) {
	return kycs.VerifyIDBySerialNumber(input)
}

func (r *mutationResolver) VerifyEmail(ctx context.Context, input model.VerificationInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetPin(ctx context.Context, input model.SetPinInput) (bool, error) {
	return users.SetUserPin(input)
}

func (r *mutationResolver) FirstTimeLogin(ctx context.Context, input model.UserLoginInput) (*model.RegistertionPayload, error) {
	return users.LoginUserByPhone(input)
}

func (r *mutationResolver) PinLogin(ctx context.Context, input model.PinLoginInput) (*model.RegistertionPayload, error) {
	return users.LoginUserByPin(input)
}

func (r *mutationResolver) ApplyLoan(ctx context.Context, input model.ApplyLoanInput) (*model.Loan, error) {
	return loans.ApplyFORLoan(input)
}

func (r *mutationResolver) SubmitKyc(ctx context.Context, input model.KycDetails) (bool, error) {
	return kycs.SaveUserKyc(input)
}

func (r *mutationResolver) VerifyKyc(ctx context.Context, input model.VerifyKycInput) (bool, error) {
	return kycs.ApproveUserKYC(input)
}

func (r *mutationResolver) RepayLoan(ctx context.Context, input model.LoanPaymentInput) (bool, error) {
	return loans.PayLoan(input)
}

func (r *mutationResolver) ChangePassword(ctx context.Context, input model.ChangePasswordInput) (bool, error) {
	return users.ChangePassword(input)
}

func (r *mutationResolver) ChangePin(ctx context.Context, input model.ChangePinInput) (bool, error) {
	return users.ChangePin(input)
}

func (r *mutationResolver) AgentLogin(ctx context.Context, input model.LoginInput) (*model.AgentAuthPayload, error) {
	return agents.LoginAgent(input)
}

func (r *queryResolver) FetchKyc(ctx context.Context, input model.FetchKycInput) ([]*model.KycQuestion, error) {
	return kycs.FetchKycs(input)
}

func (r *queryResolver) FetchUserDetails(ctx context.Context, input model.FetchUserInput) (*model.User, error) {
	return users.FetchCustomerDetails(input)
}

func (r *queryResolver) FetchLoanDetails(ctx context.Context, input model.FetchLoanInput) (*model.Loan, error) {
	return loans.FetchLoanDetails(input)
}

func (r *queryResolver) FetchAllUserLoans(ctx context.Context, input model.FetchUserLoansInput) ([]*model.Loan, error) {
	return loans.FetchUserLoans(input)
}

func (r *queryResolver) FetchUsers(ctx context.Context, input *model.PaginationInput) ([]*model.User, error) {
	return users.FetchAllCustomers(input)
}

func (r *queryResolver) FetchAllLoans(ctx context.Context, input model.LoansInput) ([]*model.Loan, error) {
	return loans.FetchLoans(input)
}

func (r *queryResolver) FetchAgents(ctx context.Context, input *model.PaginationInput) ([]*model.Agent, error) {
	return agents.FetchAllAgents(input)
}

func (r *queryResolver) FetchAgentDetails(ctx context.Context, input model.FetchAgentInput) (*model.Agent, error) {
	return agents.FetchAgentProfile(input)
}

func (r *queryResolver) SearchUser(ctx context.Context, input string) (string, error) {
	return users.SearchUserPhone(input)
}

func (r *queryResolver) FetchFaqs(ctx context.Context, input *model.PaginationInput) ([]*model.Faq, error) {
	var AllFaqs = []*model.Faq{
		{Question: "What is Ikredit Paybill number?", Answer: "Our paybill number is "},
		{Question: "Can I use another number to repay my loan?", Answer: "Yes, As long as you input the registered number as the account number"},
		{Question: "What penalty fee are applied if I don't repay my loan by due date?", Answer: "A penalty of 2% will be applied on the outstanding balance each day of delay. To avoid being penalized, repay on time and build your credit score."},
		{Question: "Are there any automatic roll-over fees?", Answer: "No.Icredit does not charge automatic roll-over fees."},
		{Question: "Will I be informed of my upcoming due date", Answer: "You will be informed by sms of your upcoming duedate"},
		{Question: "Can I be allowed to change my details after registration?", Answer: "Your user info details cannot be changed to maintain the authenticity of the info provided. However, some details  like the pin/password you will be allowed to changed."},
		{Question: "What data is used by ikreditt during its decision process.", Answer: "We take many factors into account when evaluating your creditworthiness including your prior history, data from your phone and information from CRB."},
		{Question: "How much can I borrow?", Answer: "We offer loans from ksh1000- ksh 50000"},
		{Question: "How can i increase my limit?", Answer: "Building a positive credit by repaying your loan on time"},
		{Question: "What repayment period do you offer?", Answer: "We offer variety of repayment periods depending on the limit taken"},
	}
	return AllFaqs, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
