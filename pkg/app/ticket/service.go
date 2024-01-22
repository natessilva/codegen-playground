package ticket

import (
	"codegen/app/pkg/app"
	"codegen/app/pkg/authn"
	"context"
	"database/sql"
	"strings"

	"codegen/app/db/model"

	"github.com/pkg/errors"
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

func (s *Service) Create(ctx context.Context, i app.Ticket) (app.ID, error) {
	user := authn.UserFromFromContext(ctx)
	id, err := s.q.CreateTicket(ctx, model.CreateTicketParams{
		SpaceID: user.SpaceID,
		Body:    i.Body,
		Subject: i.Subject,
	})
	if err != nil {
		return app.ID{}, errors.Wrap(err, "query error")
	}
	return app.ID{
		ID: id,
	}, nil
}

func (s *Service) Get(ctx context.Context, i app.ID) (app.GetTicketResponse, error) {
	user := authn.UserFromFromContext(ctx)
	ticket, err := s.q.GetTicket(ctx, model.GetTicketParams{
		SpaceID: user.SpaceID,
		ID:      i.ID,
	})
	if err == sql.ErrNoRows {
		return app.GetTicketResponse{
			OK: false,
		}, nil
	}
	if err != nil {
		return app.GetTicketResponse{}, errors.Wrap(err, "sql error")
	}
	return app.GetTicketResponse{
		Ticket: app.Ticket{
			Body:    ticket.Body,
			Subject: ticket.Subject,
		},
		OK: true,
	}, nil
}

func (s *Service) Update(ctx context.Context, i app.UpdateTicketInput) (app.OK, error) {
	user := authn.UserFromFromContext(ctx)
	result, err := s.q.UpdateTicket(ctx, model.UpdateTicketParams(model.UpdateTicketParams{
		SpaceID: user.SpaceID,
		ID:      i.ID,
		Subject: i.Subject,
		Body:    i.Body,
	}))
	if err != nil {
		return app.OK{}, errors.Wrap(err, "query error")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return app.OK{}, errors.Wrap(err, "rows affected error")
	}
	return app.OK{
		OK: rows == 1,
	}, nil
}

func (s *Service) Assign(ctx context.Context, i app.AssignInput) (app.OK, error) {
	user := authn.UserFromFromContext(ctx)
	err := s.q.Assign(ctx, model.AssignParams{
		SpaceID:    user.SpaceID,
		TicketID:   i.TicketID,
		IdentityID: i.UserID,
	})
	if err != nil {
		// a FK constraint violation here means either the ticket doesn't exist
		// or the user doesn't.
		if strings.Contains(err.Error(), `insert or update on table "assignee" violates foreign key constraint`) {
			return app.OK{
				OK: false,
			}, nil
		}
		return app.OK{}, errors.Wrap(err, "sql error")
	}
	return app.OK{
		OK: true,
	}, nil
}

func (s *Service) AssignSelf(ctx context.Context, i app.ID) (app.OK, error) {
	user := authn.UserFromFromContext(ctx)
	return s.Assign(ctx, app.AssignInput{
		TicketID: i.ID,
		UserID:   user.ID,
	})
}
