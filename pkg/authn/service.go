package authn

import (
	"codegen/app/db/model"
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const BcryptCost = 10

type Service struct {
	signingKey string
	q          *model.Queries
	db         *sql.DB
}

func NewService(q *model.Queries, db *sql.DB, signingKey string) *Service {
	return &Service{
		signingKey: signingKey,
		q:          q,
		db:         db,
	}
}

func (s *Service) Signup(ctx context.Context, i AuthInput) (AuthOutput, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(i.Password), BcryptCost)
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "hashing error")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "tx begin")
	}
	defer tx.Rollback()
	q := s.q.WithTx(tx)

	spaceID, err := q.CreateSpace(ctx, "")
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "create Space")
	}
	id, err := q.CreateIdentity(ctx, model.CreateIdentityParams{
		Email:          i.Email,
		PasswordHash:   hash,
		CurrentSpaceID: spaceID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return AuthOutput{
				OK: false,
			}, nil
		}
		return AuthOutput{}, errors.Wrap(err, "create user")
	}
	err = q.CreateUser(ctx, model.CreateUserParams{
		SpaceID:    spaceID,
		IdentityID: id,
	})
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "attach user")
	}
	t, err := GetTokenForUser(s.signingKey, spaceID, id, time.Hour*2)
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "signtoken")
	}
	err = tx.Commit()
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "commit")
	}
	return AuthOutput{
		Token: t,
		OK:    true,
	}, nil
}

func (s *Service) Login(ctx context.Context, i AuthInput) (AuthOutput, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "begin")
	}
	defer tx.Rollback()
	q := s.q.WithTx(tx)
	identity, err := q.GetIdentityByEmail(ctx, i.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return AuthOutput{
				OK: false,
			}, nil
		}
		return AuthOutput{}, errors.Wrap(err, "query error")
	}
	err = bcrypt.CompareHashAndPassword(identity.PasswordHash, []byte(i.Password))
	if err != nil {
		return AuthOutput{
			OK: false,
		}, nil
	}
	user, err := q.GetUser(ctx, model.GetUserParams{
		SpaceID:    identity.CurrentSpaceID.Int32,
		IdentityID: identity.ID,
	})
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "get Space user")
	}
	t, err := GetTokenForUser(s.signingKey, user.SpaceID, user.IdentityID, time.Hour*24*30)
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
