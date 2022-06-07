package ticket

import (
	"codegen/app/pkg/app"
	"codegen/app/pkg/authn"
	"context"
	"database/sql"

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

func (s *Service) Create(ctx context.Context, i app.NewTicket) (app.ID, error) {
	identity := authn.IdentityFromFromContext(ctx)
	id, err := s.q.CreateTicket(ctx, model.CreateTicketParams{
		WorkspaceID: int32(identity.WorkspaceID),
		Subject:     i.Subject,
		Description: i.Description,
	})
	if err != nil {
		return app.ID{}, errors.Wrap(err, "create")
	}
	return app.ID{ID: int(id)}, nil
}

func (s *Service) Get(ctx context.Context, i app.ID) (app.GetTicket, error) {
	identity := authn.IdentityFromFromContext(ctx)
	ticket, err := s.q.GetTicket(ctx, int32(i.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			return app.GetTicket{OK: false}, nil
		}
		return app.GetTicket{}, errors.Wrap(err, "get")
	}
	if identity.WorkspaceID != int(ticket.WorkspaceID) {
		return app.GetTicket{
			OK: false,
		}, nil
	}
	return app.GetTicket{
		Ticket: toAppTicket(ticket),
		OK:     true,
	}, nil
}

func (s *Service) Update(ctx context.Context, i app.Ticket) (app.OK, error) {
	identity := authn.IdentityFromFromContext(ctx)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return app.OK{}, errors.Wrap(err, "begin")
	}
	defer tx.Rollback()

	q := s.q.WithTx(tx)
	ticket, err := q.GetTicket(ctx, int32(i.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			return app.OK{OK: false}, nil
		}
		return app.OK{}, errors.Wrap(err, "get")
	}
	if identity.WorkspaceID != int(ticket.WorkspaceID) {
		return app.OK{OK: false}, nil
	}
	err = q.UpdateTicket(ctx, model.UpdateTicketParams{
		ID:          int32(i.ID),
		Subject:     i.Subject,
		Description: i.Description,
		Status:      model.TicketStatus(i.Status),
	})
	if err != nil {
		return app.OK{}, errors.Wrap(err, "update")
	}
	err = tx.Commit()
	if err != nil {
		return app.OK{}, errors.Wrap(err, "commit")
	}
	return app.OK{OK: true}, nil
}

func (s *Service) List(ctx context.Context, i app.Empty) (app.Tickets, error) {
	identity := authn.IdentityFromFromContext(ctx)
	list, err := s.q.GetWorkspaceTickets(ctx, int32(identity.WorkspaceID))
	if err != nil {
		return app.Tickets{}, errors.Wrap(err, "get")
	}
	tickets := make([]app.Ticket, len(list))
	for i := 0; i < len(list); i++ {
		tickets[i] = toAppTicket(list[i])
	}
	return app.Tickets{
		List: tickets,
	}, nil
}

func toAppTicket(ticket model.Ticket) app.Ticket {
	return app.Ticket{
		ID:          int(ticket.ID),
		Subject:     ticket.Subject,
		Description: ticket.Description,
		Status:      string(ticket.Status),
	}
}
