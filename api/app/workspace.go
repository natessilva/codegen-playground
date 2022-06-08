package app

type WorkspaceService interface {
	Get(Empty) Workspace
	Update(Workspace) Empty
	Create(Workspace) ID
	List(Empty) Workspaces
	Switch(ID) AuthOutput
	AddUser(AddUserInput) OK
}

type Workspace struct {
	Name string
}

type WorkspaceListItem struct {
	Name string
	ID   int
}

type Workspaces struct {
	List []WorkspaceListItem
}

type AuthOutput struct {
	Token string
	OK    bool
}

type AddUserInput struct {
	Email string
}
