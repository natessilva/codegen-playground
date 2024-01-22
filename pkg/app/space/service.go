package space

import (
	"codegen/app/pkg/app"
	"codegen/app/pkg/authn"
	"context"
	"database/sql"

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

func (s *Service) Get(ctx context.Context, i app.Empty) (app.Space, error) {
	user := authn.UserFromFromContext(ctx)
	space, err := s.q.GetSpace(ctx, int32(user.SpaceID))
	if err != nil {
		return app.Space{}, errors.Wrap(err, "get")
	}
	return app.Space{
		Name: space.Name,
	}, nil
}

func (s *Service) Update(ctx context.Context, i app.Space) (app.Empty, error) {
	user := authn.UserFromFromContext(ctx)
	err := s.q.UpdateSpace(ctx, model.UpdateSpaceParams{
		ID:   int32(user.SpaceID),
		Name: i.Name,
	})
	return app.Empty{}, err
}

func (s *Service) GetUsers(ctx context.Context, i app.Empty) (app.GetUsersResponse, error) {
	user := authn.UserFromFromContext(ctx)
	dbUsers, err := s.q.GetUsersBySpace(ctx, user.SpaceID)
	if err != nil {
		return app.GetUsersResponse{}, errors.Wrap(err, "sql error")
	}
	users := make([]app.User, 0, len(dbUsers))
	for _, user := range dbUsers {
		users = append(users, app.User{
			ID:   user.ID,
			Name: user.Name,
		})
	}
	return app.GetUsersResponse{
		Users: users,
	}, nil
}

func (s *Service) List(ctx context.Context, i app.Empty) (app.ListSpacesResponse, error) {
	user := authn.UserFromFromContext(ctx)
	spaces, err := s.q.ListSpaces(ctx, user.ID)
	if err != nil {
		return app.ListSpacesResponse{}, err
	}
	return app.ListSpacesResponse{
		Spaces: ToAppSpace(spaces),
	}, nil
}

func ToAppSpace(modelSpaces []model.Space) []app.Space {
	appSpaces := make([]app.Space, 0, len(modelSpaces))
	for _, s := range modelSpaces {
		appSpaces = append(appSpaces, app.Space{
			Name: s.Name,
		})
	}
	return appSpaces
}
