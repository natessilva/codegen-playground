-- name: CreateWorkspace :one
insert into workspace(name)
values($1)
returning id;

-- name: GetWorkspace :one
select
  *
from workspace
where id = $1;

-- name: UpdateWorkspace :exec
update workspace
set name = $2
where id = $1;

-- name: GetUserWorkspaces :many
select
    w.*
from workspace_user wu
join workspace w on wu.workspace_id = w.id
where wu.user_id = $1
order by w.id;