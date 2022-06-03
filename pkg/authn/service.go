package authn

import (
	"codegen/app/db/model/userdb"
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const BcryptCost = 10

type Service struct {
	signingKey string
	user       *userdb.Queries
}

func NewService(user *userdb.Queries, signingKey string) *Service {
	return &Service{
		user:       user,
		signingKey: signingKey,
	}
}

func (s *Service) Signup(ctx context.Context, i AuthInput) (AuthOutput, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(i.Password), BcryptCost)
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "hashing error")
	}
	userId, err := s.user.Create(ctx, userdb.CreateParams{
		Email:        i.Email,
		PasswordHash: hash,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return AuthOutput{
				OK: false,
			}, nil
		}
		return AuthOutput{}, errors.Wrap(err, "query error")
	}
	t, err := getTokenForUser(s.signingKey, int(userId), time.Hour*2)
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "signing token")
	}
	return AuthOutput{
		Token: t,
		OK:    true,
	}, nil
}

func (s *Service) Login(ctx context.Context, i AuthInput) (AuthOutput, error) {
	user, err := s.user.GetByEmail(ctx, i.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return AuthOutput{
				OK: false,
			}, nil
		}
		return AuthOutput{}, errors.Wrap(err, "query error")
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(i.Password))
	if err != nil {
		return AuthOutput{
			OK: false,
		}, nil
	}
	t, err := getTokenForUser(s.signingKey, int(user.ID), time.Hour*24*30)
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "signing token")
	}
	return AuthOutput{
		Token: t,
		OK:    true,
	}, nil
}

func (s *Service) ExchangeEmailToken(context.Context, ExchangeEmailTokenInput) (AuthOutput, error) {
	return AuthOutput{}, nil
}

func (s *Service) ResetPassword(context.Context, ResetPasswordInput) (ResetPasswordOutput, error) {
	return ResetPasswordOutput{}, nil
}
