package app

type SpaceService interface {
	Get(Empty) Space
	Update(Space) Empty
	GetUsers(Empty) GetUsersResponse
	List(Empty) ListSpacesResponse
}

type Space struct {
	Name string
}

type ListSpacesResponse struct {
	Spaces []Space
}

type GetUsersResponse struct {
	Users []User
}

type ID struct {
	ID int32
}
