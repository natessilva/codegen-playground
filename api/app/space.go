package app

type SpaceService interface {
	Get(Empty) Space
	Update(Space) Empty
	GetUsers(Empty) GetUsersResponse
}

type Space struct {
	Name string
}

type GetUsersResponse struct {
	Users []User
}

type ID struct {
	ID int32
}
