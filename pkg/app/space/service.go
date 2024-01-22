package space

import (
	"codegen/app/pkg/app"
	"codegen/app/pkg/authn"
	"codegen/app/pkg/slices"
	"context"

	"codegen/app/db/model"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	q          *model.Queries
	db         *pgxpool.Pool
	signingKey string
}

func NewService(q *model.Queries, db *pgxpool.Pool, signingKey string) *Service {
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
		Spaces: slices.Map(spaces, func(s model.Space) app.Space {
			return app.Space{
				Name: s.Name,
			}
		}),
	}, nil
}
