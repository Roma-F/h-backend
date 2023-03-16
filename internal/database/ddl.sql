create table if not exists ads
(
    id         bigint unsigned not null primary key auto_increment,
    user_id    bigint unsigned not null,
    created_at timestamp default now(),
    address_status  tinyint default 0,
    rooms_status    tinyint default 0,
    bathroom_status tinyint default 0,

    key (user_id)
);

 create table ads_opts_types
 (
     opt_type   bigint unsigned not null primary key,
     name       varchar(255)    not null,
     value_type int unsigned    not null
 );

create table if not exists ads_opts
(
    id             bigint unsigned not null primary key auto_increment,
    ad_id          bigint unsigned not null,
    opt_type       bigint unsigned not null,

    opt_value_str  varchar(255),
    opt_value_uint bigint unsigned,
    opt_value_bool boolean,
    opt_value_blob mediumblob,

    unique key (ad_id, opt_type)
);

create table cities
(
    id         bigint unsigned not null primary key auto_increment,
    name       varchar(255)    not null,
    translit   varchar(255)    not null,
    created_at timestamp       default now()
);

create index ad_id_2
    on ads_opts (ad_id, opt_type, opt_value_str);

create index ad_val_bool
    on ads_opts (opt_type, opt_value_bool);

create index ad_val_str
    on ads_opts (opt_type, opt_value_str);

create index ad_val_uint
    on ads_opts (opt_type, opt_value_uint);

-- create table if not exists users
-- (
--     id         bigint unsigned not null primary key auto_increment,
--     created_at timestamp default now()
-- );

-- create table if not exists  users_opts
-- (
--     id             bigint unsigned not null primary key auto_increment,
--     user_id        bigint unsigned not null,
--     opt_type       bigint unsigned not null,

--     opt_value_str  varchar(255),
--     opt_value_uint bigint unsigned,
--     opt_value_bool boolean,
--     opt_value_blob mediumblob,

--     unique key (user_id, opt_type),
--     key (user_id, opt_type, opt_value_str)
-- );

-- create table if not exists token
-- (
--     uuid       varchar(36) primary key ,
--     user_id    bigint      not null,
--     active     boolean     default false,
--     expiry     timestamp   not null,
--     created_at timestamp   default now(),
--     data       text,
--     key (user_id, active)
-- );

-- ALTER TABLE ads ADD bathroom_status tinyint default 0;