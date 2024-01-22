package user

import (
	"codegen/app/pkg/app"
	"codegen/app/pkg/authn"
	"context"

	"codegen/app/db/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	q  *model.Queries
	db *pgxpool.Pool
}

func NewService(q *model.Queries, db *pgxpool.Pool) *Service {
	return &Service{
		q:  q,
		db: db,
	}
}

func (s *Service) Get(ctx context.Context, i app.Empty) (app.User, error) {
	identity := authn.UserFromFromContext(ctx)
	user, err := s.q.GetIdentity(ctx, identity.ID)
	if err != nil {
		return app.User{}, errors.Wrap(err, "query error")
	}
	return app.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

func (s *Service) SetPassword(ctx context.Context, i app.SetPasswordInput) (app.OK, error) {
	identity := authn.UserFromFromContext(ctx)
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return app.OK{}, errors.Wrap(err, "begin error")
	}
	defer tx.Rollback(ctx)
	q := s.q.WithTx(tx)
	u, err := q.GetIdentity(ctx, identity.ID)
	if err != nil {
		return app.OK{}, errors.Wrap(err, "query error getting hash")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(i.OldPassword))
	if err != nil {
		return app.OK{
			OK: false,
		}, nil
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(i.NewPassword), authn.BcryptCost)
	if err != nil {
		return app.OK{}, errors.Wrap(err, "hashing error")
	}
	err = q.SetIdentityPassword(ctx, model.SetIdentityPasswordParams{
		PasswordHash: newHash,
		ID:           identity.ID,
	})
	if err != nil {
		return app.OK{}, errors.Wrap(err, "query error setting hash")
	}
	err = tx.Commit(ctx)
	if err != nil {
		return app.OK{}, errors.Wrap(err, "commit erorr")
	}
	return app.OK{
		OK: true,
	}, nil
}

func (s *Service) Update(ctx context.Context, i app.UpdateUserInput) (app.Empty, error) {
	identity := authn.UserFromFromContext(ctx)
	err := s.q.UpdateIdentity(ctx, model.UpdateIdentityParams{
		ID:   identity.ID,
		Name: i.Name,
	})
	if err != nil {
		return app.Empty{}, errors.Wrap(err, "query error")
	}
	return app.Empty{}, nil
}
