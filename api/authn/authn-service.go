package authn

type AuthNService interface {
	Signup(AuthInput) AuthOutput
	Login(AuthInput) AuthOutput
	ResetPassword(ResetPasswordInput) ResetPasswordOutput
	ExchangeEmailToken(ExchangeEmailTokenInput) AuthOutput
}

type AuthInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	Token string
	OK    bool
}

type ResetPasswordInput struct {
	Email string
}

type ResetPasswordOutput struct {
	OK bool
}

type ExchangeEmailTokenInput struct {
	Token string
}
