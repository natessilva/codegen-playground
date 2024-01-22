// Code generated - DO NOT EDIT.

package authn

import (
	"context"
	"net/http"

	"codegen/app/pkg/apimux"
)

type AuthNService interface {
	ExchangeEmailToken(context.Context, ExchangeEmailTokenInput) (AuthOutput, error)
	Login(context.Context, AuthInput) (AuthOutput, error)
	ResetPassword(context.Context, ResetPasswordInput) (ResetPasswordOutput, error)
	Signup(context.Context, AuthInput) (AuthOutput, error)
}

type authNServiceServer struct {
	authNService AuthNService
}

// Register the implementation of AuthNService with the apimux Server.
func RegisterAuthNService(server *apimux.Server, authNService AuthNService) {
	handler := &authNServiceServer{
		authNService: authNService,
	}
	
	server.Register("AuthNService", "ExchangeEmailToken", handler.handleExchangeEmailToken)
	server.Register("AuthNService", "Login", handler.handleLogin)
	server.Register("AuthNService", "ResetPassword", handler.handleResetPassword)
	server.Register("AuthNService", "Signup", handler.handleSignup)
}

func (s *authNServiceServer) handleExchangeEmailToken(w http.ResponseWriter, r *http.Request) {
	var input ExchangeEmailTokenInput
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.authNService.ExchangeEmailToken(r.Context(), input)
	})
}

func (s *authNServiceServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var input AuthInput
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.authNService.Login(r.Context(), input)
	})
}

func (s *authNServiceServer) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var input ResetPasswordInput
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.authNService.ResetPassword(r.Context(), input)
	})
}

func (s *authNServiceServer) handleSignup(w http.ResponseWriter, r *http.Request) {
	var input AuthInput
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.authNService.Signup(r.Context(), input)
	})
}

type ResetPasswordOutput struct {
	OK bool `json:"ok"`
}

type ExchangeEmailTokenInput struct {
	Token string `json:"token"`
}

type AuthOutput struct {
	Token string `json:"token"`
	OK bool `json:"ok"`
}

type AuthInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordInput struct {
	Email string `json:"email"`
}

