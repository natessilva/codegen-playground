-- name: CreateUser :one
insert into "user"(email, password_hash, current_workspace_id)
values($1,$2, $3)
on conflict do nothing
returning id;

-- name: GetUser :one
select
  *
from "user"
where id = $1;

-- name: GetUserByEmail :one
select
  *
from "user"
where email = $1;

-- name: UpdateUser :exec
update "user"
set name = $2
where id = $1;

-- name: SetUserPassword :exec
update "user"
set password_hash = $2
where id = $1;

-- name: SetUserWorkspace :exec
update "user"
set current_workspace_id = $2
where id = $1;