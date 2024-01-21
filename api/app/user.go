package app

type UserService interface {
	Get(Empty) User
	Update(UpdateUserInput) Empty
	SetPassword(SetPasswordInput) OK
}

type User struct {
	ID   int32
	Name string
}

type UpdateUserInput struct {
	Name string
}
type SetPasswordInput struct {
	OldPassword string
	NewPassword string
}
