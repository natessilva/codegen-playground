package authn

import (
	"codegen/app/db/model"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const BcryptCost = 10

type Service struct {
	signingKey string
	q          *model.Queries
	db         *pgxpool.Pool
}

func NewService(q *model.Queries, db *pgxpool.Pool, signingKey string) *Service {
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

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "tx begin")
	}
	defer tx.Rollback(ctx)
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
		if err == pgx.ErrNoRows {
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
	err = tx.Commit(ctx)
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "commit")
	}
	return AuthOutput{
		Token: t,
		OK:    true,
	}, nil
}

func (s *Service) Login(ctx context.Context, i AuthInput) (AuthOutput, error) {
	identity, err := s.q.GetIdentityByEmail(ctx, i.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
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
	user, err := s.q.GetUser(ctx, model.GetUserParams{
		SpaceID:    identity.CurrentSpaceID.Int32,
		IdentityID: identity.ID,
	})
	if err != nil {
		return AuthOutput{}, errors.Wrap(err, "get user")
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
