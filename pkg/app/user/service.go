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
	model *model.Queries
	db    *sql.DB
}

func NewService(model *model.Queries, db *sql.DB) *Service {
	return &Service{
		model: model,
		db:    db,
	}
}

func (s *Service) Get(ctx context.Context, i app.Empty) (app.UserInfo, error) {
	identity := authn.IdentityFromFromContext(ctx)
	name, err := s.model.GetUser(ctx, int32(identity.UserID))
	if err != nil {
		return app.UserInfo{}, errors.Wrap(err, "query error")
	}
	return app.UserInfo{
		Name: name,
	}, nil
}

func (s *Service) SetPassword(ctx context.Context, i app.SetPasswordInput) (app.OK, error) {
	identity := authn.IdentityFromFromContext(ctx)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return app.OK{}, errors.Wrap(err, "begin error")
	}
	defer tx.Rollback()
	txModel := s.model.WithTx(tx)
	oldHash, err := txModel.GetPasswordById(ctx, int32(identity.UserID))
	if err != nil {
		return app.OK{}, errors.Wrap(err, "query error getting hash")
	}
	err = bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(i.OldPassword))
	if err != nil {
		return app.OK{
			OK: false,
		}, nil
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(i.NewPassword), authn.BcryptCost)
	if err != nil {
		return app.OK{}, errors.Wrap(err, "hashing error")
	}
	err = txModel.SetPassword(ctx, model.SetPasswordParams{
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
	err := s.model.UpdateUser(ctx, model.UpdateUserParams{
		ID:   int32(identity.UserID),
		Name: i.Name,
	})
	if err != nil {
		return app.Empty{}, errors.Wrap(err, "query error")
	}
	return app.Empty{}, nil
}
