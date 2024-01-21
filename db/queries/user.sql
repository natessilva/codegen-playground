-- name: CreateUser :exec
insert into "user"(space_id, identity_id)
values ($1,$2);

-- name: GetUser :one
select *
from "user"
where space_id = $1 and identity_id = $2;

-- name: GetUsersBySpace :many
select i.*
from "user" u
join identity i on u.identity_id = i.id
where space_id = $1;