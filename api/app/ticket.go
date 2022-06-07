package app

type TicketService interface {
	Create(NewTicket) ID
	Get(ID) GetTicket
	Update(Ticket) OK
	List(Empty) Tickets
}

type TicketStatus string

const (
	TicketStatusOpen     TicketStatus = "open"
	TicketStatusClosed   TicketStatus = "closed"
	TicketStatusArchived TicketStatus = "archived"
)

type NewTicket struct {
	Subject     string
	Description string
}

type Ticket struct {
	ID          int
	Subject     string
	Description string
	Status      TicketStatus
}

type GetTicket struct {
	Ticket Ticket
	OK     bool
}

type Tickets struct {
	List []Ticket
}
