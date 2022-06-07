-- migrate:up
create table workspace(
    id serial primary key,
    name text not null default ''
);

create table workspace_user(
    id serial primary key,
    workspace_id int not null references workspace (id),
    user_id int not null references "user" (id),
    unique(workspace_id,user_id)
);

alter table "user" add column current_workspace_id int not null references workspace (id);

-- migrate:down
alter table "user" drop column current_workspace_id;
drop table workspace_user;
drop table workspace;