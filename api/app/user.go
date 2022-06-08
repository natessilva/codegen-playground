package app

type UserService interface {
	Get(Empty) User
	Update(User) Empty
	SetPassword(SetPasswordInput) OK
}

type User struct {
	Name string
}
type SetPasswordInput struct {
	OldPassword string
	NewPassword string
}
