DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT FROM pg_database WHERE datname = 'order_service'
        ) THEN
            EXECUTE 'CREATE DATABASE order_service';
        END IF;
    END $$;create table if not exists orders
(
    id    serial not null
        constraint orders_pk
            primary key,
    item  TEXT   not null,
    quantity integer
);