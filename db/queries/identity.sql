-- name: CreateIdentity :one
insert into identity(email, password_hash, name, current_space_id)
values($1,$2, $3, @current_space_id::int)
on conflict do nothing
returning id;

-- name: GetIdentity :one
select
  *
from identity
where id = $1;

-- name: GetIdentityByEmail :one
select
  *
from identity
where email = $1;

-- name: UpdateIdentity :exec
update identity
set name = $2
where id = $1;

-- name: SetIdentityPassword :exec
update identity
set password_hash = $2
where id = $1;

-- name: SetIdentityCurrentSpace :exec
update identity
set current_space_id = $2
where id = $1;