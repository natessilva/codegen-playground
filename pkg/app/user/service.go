package user

import (
	"codegen/app/pkg/app"
	"codegen/app/pkg/authn"
	"context"
	"database/sql"

	"codegen/app/db/model"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	q  *model.Queries
	db *sql.DB
}

func NewService(q *model.Queries, db *sql.DB) *Service {
	return &Service{
		q:  q,
		db: db,
	}
}

func (s *Service) Get(ctx context.Context, i app.Empty) (app.User, error) {
	identity := authn.IdentityFromFromContext(ctx)
	user, err := s.q.GetUser(ctx, int32(identity.UserID))
	if err != nil {
		return app.User{}, errors.Wrap(err, "query error")
	}
	return app.User{
		Name: user.Name,
	}, nil
}

func (s *Service) SetPassword(ctx context.Context, i app.SetPasswordInput) (app.OK, error) {
	identity := authn.IdentityFromFromContext(ctx)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return app.OK{}, errors.Wrap(err, "begin error")
	}
	defer tx.Rollback()
	q := s.q.WithTx(tx)
	u, err := q.GetUser(ctx, int32(identity.UserID))
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
	err = q.SetUserPassword(ctx, model.SetUserPasswordParams{
		PasswordHash: newHash,
		ID:           int32(identity.UserID),
	})
	if err != nil {
		return app.OK{}, errors.Wrap(err, "query error setting hash")
	}
	err = tx.Commit()
	if err != nil {
		return app.OK{}, errors.Wrap(err, "commit erorr")
	}
	return app.OK{
		OK: true,
	}, nil
}

func (s *Service) Update(ctx context.Context, i app.User) (app.Empty, error) {
	identity := authn.IdentityFromFromContext(ctx)
	err := s.q.UpdateUser(ctx, model.UpdateUserParams{
		ID:   int32(identity.UserID),
		Name: i.Name,
	})
	if err != nil {
		return app.Empty{}, errors.Wrap(err, "query error")
	}
	return app.Empty{}, nil
}
