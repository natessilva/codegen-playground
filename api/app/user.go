package app

type UserService interface {
	Get(Empty) UserInfo
	Update(UserInfo) Empty
	SetPassword(SetPasswordInput) OK
}

type UserInfo struct {
	Name string
}
type SetPasswordInput struct {
	OldPassword string
	NewPassword string
}
