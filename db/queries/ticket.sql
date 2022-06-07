-- name: CreateTicket :one
insert into ticket(workspace_id, subject, description)
values($1,$2, $3)
returning id;

-- name: GetTicket :one
select
  *
from ticket
where id = $1;

-- name: GetWorkspaceTickets :many
select
  *
from ticket
where workspace_id = $1;

-- name: UpdateTicket :exec
update ticket
set subject = $2,
description = $3,
status = $4
where id = $1;