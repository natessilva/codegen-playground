package user

import (
	"codegen/app/pkg/app"
	"codegen/app/pkg/authn"
	"context"
	"database/sql"

	"codegen/app/db/model/userdb"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	user *userdb.Queries
	db   *sql.DB
}

func NewService(user *userdb.Queries, db *sql.DB) *Service {
	return &Service{
		user: user,
		db:   db,
	}
}

func (s *Service) Get(ctx context.Context, i app.Empty) (app.UserInfo, error) {
	identity := authn.IdentityFromFromContext(ctx)
	user, err := s.user.Get(ctx, int32(identity.UserID))
	if err != nil {
		return app.UserInfo{}, errors.Wrap(err, "query error")
	}
	return app.UserInfo{
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
	userTx := s.user.WithTx(tx)
	u, err := userTx.Get(ctx, int32(identity.UserID))
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
	err = userTx.SetPassword(ctx, userdb.SetPasswordParams{
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

func (s *Service) Update(ctx context.Context, i app.UserInfo) (app.Empty, error) {
	identity := authn.IdentityFromFromContext(ctx)
	err := s.user.Update(ctx, userdb.UpdateParams{
		ID:   int32(identity.UserID),
		Name: i.Name,
	})
	if err != nil {
		return app.Empty{}, errors.Wrap(err, "query error")
	}
	return app.Empty{}, nil
}
