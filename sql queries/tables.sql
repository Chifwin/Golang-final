-- create database go_final;

create extension pgcrypto;

drop type if exists user_role cascade;
create type user_role as Enum (
    'admin',
    'seller',
    'buyer'
    );

drop table if exists users cascade;
create table users
(
    id       serial primary key,
    username varchar(50)  not null unique,
    password text         not null,
    name     varchar(100) not null,
    role     user_role
);
create or replace function restrict_user_role_change() returns trigger as
$$
begin
    if (new.role != old.role) then
        raise exception 'users can not change roles';
    end if;
    return new;
end;
$$
    language plpgsql;
create or replace trigger restrict_user_role_change_tg
    before update
    on users
    for each row
execute function restrict_user_role_change();

drop table if exists products cascade;
create table products
(
    id          serial primary key,
    name        varchar(255) not null,
    description text
);

drop table if exists product_seller cascade;
create table product_seller
(
    product_id int references products (id) on update cascade on delete restrict,
    seller_id  int     references users (id) on update cascade on delete set null,
    quantity   int check ( quantity >= 0 ),
    cost       numeric(10, 2) check ( quantity > 0 ),
    published  boolean not null,
    primary key (product_id, seller_id)
);
create or replace function user_is_seller() returns trigger as
$$
begin
    if (select count(*) from users where id = new.seller_id and role = 'seller') != 1 then
        raise exception 'there is no seller with this id';
    end if;
    return new;
end;
$$
    language plpgsql;
create or replace trigger product_seller_tg
    before insert or update
    on product_seller
    for each row
execute function user_is_seller();

drop table if exists purchases cascade;
create table purchases
(
    id         serial primary key,
    buyer_id   int references users (id) on update cascade on delete cascade,
    product_id int references products (id) on update cascade on delete restrict,
    seller_id  int references users (id) on update cascade on delete set null,
    date       date default now(),
    quantity   int check ( quantity > 0 )
);
create or replace trigger product_seller_tg
    before insert or update
    on purchases
    for each row
execute function user_is_seller();

drop table if exists scores cascade;
create table scores
(
    purchase_id int primary key references purchases (id) on delete cascade,
    rating      numeric(3, 2) check ( rating >= 0 and rating <= 5 ),
    comment     text
);

