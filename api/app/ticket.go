package app

type TicketService interface {
	Create(NewTicket) ID
	Get(ID) GetTicket
	Update(Ticket) OK
	List(Empty) Tickets
}

type NewTicket struct {
	Subject     string
	Description string
}

type Ticket struct {
	ID          int
	Subject     string
	Description string
	Status      string
}

type GetTicket struct {
	Ticket Ticket
	OK     bool
}

type Tickets struct {
	List []Ticket
}
