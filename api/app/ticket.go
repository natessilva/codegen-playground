package app

type TicketService interface {
	Create(Ticket) ID
	Get(ID) GetTicketResponse
	Update(UpdateTicketInput) Empty
	Assign(AssignInput) OK
	AssignSelf(ID) OK
}

type Ticket struct {
	Subject string
	Body    string
}

type GetTicketResponse struct {
	Ticket Ticket
	OK     bool
}
type UpdateTicketInput struct {
	ID      int32
	Subject string
	Body    string
}

type AssignInput struct {
	TicketID int32
	UserID   int32
}
