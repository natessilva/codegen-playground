package app

type WorkspaceService interface {
	Get(Empty) WorkspaceInfo
	Update(WorkspaceInfo) Empty
	Create(WorkspaceInfo) ID
	List(Empty) Workspaces
	Switch(ID) AuthOutput
}

type WorkspaceInfo struct {
	Name string
}

type WorkspaceListItem struct {
	Name string
	ID   int
}

type ID struct {
	ID int
}

type Workspaces struct {
	List []WorkspaceListItem
}

type AuthOutput struct {
	Token string
	OK    bool
}
