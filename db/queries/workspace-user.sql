-- name: CreateWorkspaceUser :one
insert into workspace_user(workspace_id, user_id)
values($1,$2)
on conflict do nothing
returning id;

-- name: GetWorkspaceUser :one
select
  *
from workspace_user
where id = $1;

-- name: GetWorkspaceUserId :one
select
  id
from workspace_user
where workspace_id = $1 and user_id = $2;