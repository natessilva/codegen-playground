-- migrate:up
create type ticket_status as enum('open','closed','archived');

create table ticket(
    id serial primary key,
    workspace_id int not null references workspace (id),
    subject text not null default '',
    description text not null default '',
    status ticket_status not null default 'open'
);

-- migrate:down
drop table ticket;