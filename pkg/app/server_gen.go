// Code generated - DO NOT EDIT.

package app

import (
	"context"
	"net/http"

	"codegen/app/pkg/apimux"
)

type SpaceService interface {
	Get(context.Context, Empty) (Space, error)
	GetUsers(context.Context, Empty) (GetUsersResponse, error)
	List(context.Context, Empty) (ListSpacesResponse, error)
	Update(context.Context, Space) (Empty, error)
}
type TicketService interface {
	Assign(context.Context, AssignInput) (OK, error)
	AssignSelf(context.Context, ID) (OK, error)
	Create(context.Context, Ticket) (ID, error)
	Get(context.Context, ID) (GetTicketResponse, error)
	Update(context.Context, UpdateTicketInput) (OK, error)
}
type UserService interface {
	Get(context.Context, Empty) (User, error)
	SetPassword(context.Context, SetPasswordInput) (OK, error)
	Update(context.Context, UpdateUserInput) (Empty, error)
}

type spaceServiceServer struct {
	spaceService SpaceService
}

// Register the implementation of SpaceService with the apimux Server.
func RegisterSpaceService(server *apimux.Server, spaceService SpaceService) {
	handler := &spaceServiceServer{
		spaceService: spaceService,
	}
	
	server.Register("SpaceService", "Get", handler.handleGet)
	server.Register("SpaceService", "GetUsers", handler.handleGetUsers)
	server.Register("SpaceService", "List", handler.handleList)
	server.Register("SpaceService", "Update", handler.handleUpdate)
}

func (s *spaceServiceServer) handleGet(w http.ResponseWriter, r *http.Request) {
	var input Empty
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.spaceService.Get(r.Context(), input)
	})
}

func (s *spaceServiceServer) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	var input Empty
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.spaceService.GetUsers(r.Context(), input)
	})
}

func (s *spaceServiceServer) handleList(w http.ResponseWriter, r *http.Request) {
	var input Empty
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.spaceService.List(r.Context(), input)
	})
}

func (s *spaceServiceServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var input Space
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.spaceService.Update(r.Context(), input)
	})
}

type ticketServiceServer struct {
	ticketService TicketService
}

// Register the implementation of TicketService with the apimux Server.
func RegisterTicketService(server *apimux.Server, ticketService TicketService) {
	handler := &ticketServiceServer{
		ticketService: ticketService,
	}
	
	server.Register("TicketService", "Assign", handler.handleAssign)
	server.Register("TicketService", "AssignSelf", handler.handleAssignSelf)
	server.Register("TicketService", "Create", handler.handleCreate)
	server.Register("TicketService", "Get", handler.handleGet)
	server.Register("TicketService", "Update", handler.handleUpdate)
}

func (s *ticketServiceServer) handleAssign(w http.ResponseWriter, r *http.Request) {
	var input AssignInput
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.ticketService.Assign(r.Context(), input)
	})
}

func (s *ticketServiceServer) handleAssignSelf(w http.ResponseWriter, r *http.Request) {
	var input ID
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.ticketService.AssignSelf(r.Context(), input)
	})
}

func (s *ticketServiceServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	var input Ticket
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.ticketService.Create(r.Context(), input)
	})
}

func (s *ticketServiceServer) handleGet(w http.ResponseWriter, r *http.Request) {
	var input ID
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.ticketService.Get(r.Context(), input)
	})
}

func (s *ticketServiceServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var input UpdateTicketInput
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.ticketService.Update(r.Context(), input)
	})
}

type userServiceServer struct {
	userService UserService
}

// Register the implementation of UserService with the apimux Server.
func RegisterUserService(server *apimux.Server, userService UserService) {
	handler := &userServiceServer{
		userService: userService,
	}
	
	server.Register("UserService", "Get", handler.handleGet)
	server.Register("UserService", "SetPassword", handler.handleSetPassword)
	server.Register("UserService", "Update", handler.handleUpdate)
}

func (s *userServiceServer) handleGet(w http.ResponseWriter, r *http.Request) {
	var input Empty
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.userService.Get(r.Context(), input)
	})
}

func (s *userServiceServer) handleSetPassword(w http.ResponseWriter, r *http.Request) {
	var input SetPasswordInput
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.userService.SetPassword(r.Context(), input)
	})
}

func (s *userServiceServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var input UpdateUserInput
	apimux.DecodeEncode(w, r, &input, func() (interface{}, error) {
		return s.userService.Update(r.Context(), input)
	})
}

type GetTicketResponse struct {
	Ticket Ticket `json:"ticket"`
	OK bool `json:"ok"`
}

type UpdateTicketInput struct {
	ID int32 `json:"id"`
	Subject string `json:"subject"`
	Body string `json:"body"`
}

type GetUsersResponse struct {
	Users []User `json:"users"`
}

type AssignInput struct {
	TicketID int32 `json:"ticketID"`
	UserID int32 `json:"userID"`
}

type Ticket struct {
	Subject string `json:"subject"`
	Body string `json:"body"`
}

type SetPasswordInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type Empty struct {
}

type Space struct {
	Name string `json:"name"`
}

type ID struct {
	ID int32 `json:"id"`
}

type UpdateUserInput struct {
	Name string `json:"name"`
}

type User struct {
	ID int32 `json:"id"`
	Name string `json:"name"`
}

type ListSpacesResponse struct {
	Spaces []Space `json:"spaces"`
}

type OK struct {
	OK bool `json:"ok"`
}

