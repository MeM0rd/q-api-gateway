create table if not exists user_roles
(
    id    int primary key generated always as identity,
    name  varchar(30) not null,
    level int not null
);

create table if not exists users
(
    id       int primary key generated always as identity,
    surname  varchar(255) not null,
    name     varchar(255) not null,
    password varchar(255) not null,
    email    varchar(255) not null unique,
    role_id  int default 1,
    constraint fk_role_id foreign key (role_id) references user_roles (id)
);

create table if not exists quotes
(
    id int primary key generated always as identity,
    user_id int not null,
    title varchar(50) not null,
    text text not null,
    constraint fk_user_id foreign key (user_id) references users(id)
);

create table if not exists sessions
(
    id int primary key generated always as identity,
    user_id int not null,
    token varchar(50) not null unique,
    expired_at timestamp not null,
    constraint fk_user_id foreign key (user_id) references users(id)
);