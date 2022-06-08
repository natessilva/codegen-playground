package workspace

import (
	"codegen/app/pkg/app"
	"codegen/app/pkg/authn"
	"context"
	"database/sql"
	"time"

	"codegen/app/db/model"

	"github.com/pkg/errors"
)

type Service struct {
	q          *model.Queries
	db         *sql.DB
	signingKey string
}

func NewService(q *model.Queries, db *sql.DB, signingKey string) *Service {
	return &Service{
		q:          q,
		db:         db,
		signingKey: signingKey,
	}
}

func (s *Service) Get(ctx context.Context, i app.Empty) (app.Workspace, error) {
	identity := authn.IdentityFromFromContext(ctx)
	user, err := s.q.GetWorkspace(ctx, int32(identity.WorkspaceID))
	if err != nil {
		return app.Workspace{}, errors.Wrap(err, "get")
	}
	return app.Workspace{
		Name: user.Name,
	}, nil
}

func (s *Service) Update(ctx context.Context, i app.Workspace) (app.Empty, error) {
	identity := authn.IdentityFromFromContext(ctx)
	err := s.q.UpdateWorkspace(ctx, model.UpdateWorkspaceParams{
		ID:   int32(identity.WorkspaceID),
		Name: i.Name,
	})
	return app.Empty{}, err
}

func (s *Service) Create(ctx context.Context, i app.Workspace) (app.ID, error) {
	identity := authn.IdentityFromFromContext(ctx)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return app.ID{}, errors.Wrap(err, "begin")
	}
	q := s.q.WithTx(tx)

	ws, err := q.CreateWorkspace(ctx, i.Name)
	if err != nil {
		return app.ID{}, errors.Wrap(err, "create workspace")
	}
	_, err = q.CreateWorkspaceUser(ctx, model.CreateWorkspaceUserParams{
		WorkspaceID: ws,
		UserID:      int32(identity.UserID),
	})
	if err != nil {
		return app.ID{}, errors.Wrap(err, "create workspace user")
	}
	err = tx.Commit()
	if err != nil {
		return app.ID{}, errors.Wrap(err, "commit")
	}
	return app.ID{
		ID: int(ws),
	}, nil
}

func (s *Service) List(ctx context.Context, _ app.Empty) (app.Workspaces, error) {
	identity := authn.IdentityFromFromContext(ctx)
	ws, err := s.q.GetUserWorkspaces(ctx, int32(identity.UserID))
	if err != nil {
		return app.Workspaces{}, errors.Wrap(err, "get user workspaces")
	}
	workspaces := make([]app.WorkspaceListItem, len(ws))
	for i := 0; i < len(ws); i++ {
		workspaces[i] = app.WorkspaceListItem{
			Name: ws[i].Name,
			ID:   int(ws[i].ID),
		}
	}
	return app.Workspaces{
		List: workspaces,
	}, nil
}

func (s *Service) Switch(ctx context.Context, i app.ID) (app.AuthOutput, error) {
	identity := authn.IdentityFromFromContext(ctx)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return app.AuthOutput{}, errors.Wrap(err, "begin")
	}
	defer tx.Rollback()
	q := s.q.WithTx(tx)
	wsUser, err := q.GetWorkspaceUserId(ctx, model.GetWorkspaceUserIdParams{
		WorkspaceID: int32(i.ID),
		UserID:      int32(identity.UserID),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return app.AuthOutput{
				OK: false,
			}, nil
		}
		return app.AuthOutput{}, errors.Wrap(err, "get workspace user")
	}
	q.SetUserWorkspace(ctx, model.SetUserWorkspaceParams{
		ID:                 int32(identity.UserID),
		CurrentWorkspaceID: int32(i.ID),
	})
	t, err := authn.GetTokenForUser(s.signingKey, int(wsUser), time.Hour*24*30)
	if err != nil {
		return app.AuthOutput{}, errors.Wrap(err, "signing token")
	}
	err = tx.Commit()
	if err != nil {
		return app.AuthOutput{}, errors.Wrap(err, "commit")
	}
	return app.AuthOutput{
		Token: t,
		OK:    true,
	}, nil
}

func (s *Service) AddUser(ctx context.Context, i app.AddUserInput) (app.OK, error) {
	identity := authn.IdentityFromFromContext(ctx)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return app.OK{}, errors.Wrap(err, "begin")
	}
	defer tx.Rollback()
	q := s.q.WithTx(tx)

	user, err := q.GetUserByEmail(ctx, i.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return app.OK{OK: false}, nil
		}
		return app.OK{}, errors.Wrap(err, "get")
	}
	_, err = q.CreateWorkspaceUser(ctx, model.CreateWorkspaceUserParams{
		WorkspaceID: int32(identity.WorkspaceID),
		UserID:      user.ID,
	})
	if err != nil && err != sql.ErrNoRows {
		return app.OK{}, errors.Wrap(err, "create workspace user")
	}
	err = tx.Commit()
	if err != nil {
		return app.OK{}, errors.Wrap(err, "commit")
	}
	return app.OK{OK: true}, nil
}
