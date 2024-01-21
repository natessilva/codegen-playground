-- name: CreateTicket :one
insert into ticket(space_id, body, subject)
values(@space_id, @body, @subject)
returning id;

-- name: GetTicket :one
select *
from ticket
where space_id = @space_id
and id = @id;

-- name: UpdateTicket :exec
update ticket set
subject = @subject,
body = @body
where space_id = @space_id
and id = @id;