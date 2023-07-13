create table actions
(
    id           integer           not null primary key auto_increment,
    user_id      integer           not null,
    username     varchar(40)       not null,
    cmd          varchar(40)       not null,
    result       integer default 1 not null,
    execute_time int               not null,
    handler      text
);

create index actions_user_id_name_execute_time_index
    on actions (user_id asc, cmd asc, execute_time desc);

create table codes
(
    code         varchar(6)        not null unique,
    user_id      integer,
    title        varchar(40),
    attempts     integer,
    max_attempts integer default 0 not null
);

create index codes_code_max_attempts_index
    on codes (code, max_attempts);

create table `groups`
(
    id    integer not null primary key auto_increment,
    title text    not null
);

create index groups_title_index on `groups` (title);

insert into `groups`(title) values ('admin');

create table invites
(
    id      integer primary key auto_increment,
    code    text          not null unique,
    active  int default 1 not null,
    id_user integer
);

create index invites_active_code_index
    on invites (active, code);


create table membership
(
    id         integer not null primary key auto_increment,
    user_id    integer not null,
    id_group   integer not null,
    valid_till integer
);

create index membership_id_group_index
    on membership (id_group);

create index membership_id_user_index
    on membership (user_id);

create table users
(
    id      integer primary key auto_increment,
    user_id integer not null unique,
    name    TEXT    not null,
    active  INT default 1
);

create index active_idx
    on users (active);

create index users_name_index
    on users (name);

