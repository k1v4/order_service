create table if not exists orders
(
    id    serial not null
        constraint orders_pk
            primary key,
    item  TEXT   not null,
    quantity integer
);