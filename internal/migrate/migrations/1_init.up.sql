create table choices (
    choice_id serial primary key,
    name varchar(255) UNIQUE 
)

create table polls (
    poll_id serial primary key,
    closed boolean default(FALSE)
)

create table choices_to_polls (
    id serial primary key,
    hero_id int not null references choices,
    poll_id int not null references polls,
    amount int not null default(0)
)