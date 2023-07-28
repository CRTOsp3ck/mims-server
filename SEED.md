##### SEED DATA SQL STATEMENTS #####
create table agent
(
    id         serial
        primary key,
    username   varchar(255) default ''::character varying not null,
    password   varchar(255) default ''::character varying not null,
    name       varchar(255) default ''::character varying not null,
    email      varchar(255) default ''::character varying not null,
    phone      varchar(255) default ''::character varying not null,
    is_owner   boolean      default false                 not null,
    created_at timestamp                                  not null,
    updated_at timestamp                                  not null
);

alter table agent
    owner to root;

create table balance
(
    id         serial
        primary key,
    bal_cash   varchar(255) default 'sb=0&eb=0'::character varying not null,
    bal_qr     varchar(255) default 'sb=0&eb=0'::character varying not null,
    created_at timestamp                                           not null,
    updated_at timestamp                                           not null
);

alter table balance
    owner to root;

create table inventory
(
    id             serial
        primary key,
    start_item_bal varchar(255) default '1=0'::character varying not null,
    end_item_bal   varchar(255) default '1=0'::character varying not null,
    created_at     timestamp                                     not null,
    updated_at     timestamp                                     not null
);

alter table inventory
    owner to root;

create table item
(
    id         serial
        primary key,
    name       varchar(255) default ''::character varying not null,
    des        varchar(255) default ''::character varying not null,
    created_at timestamp                                  not null,
    updated_at timestamp                                  not null
);

alter table item
    owner to root;

create table operation
(
    id                 serial
        primary key,
    start_time         timestamp                                  not null,
    end_time           timestamp                                  not null,
    location           varchar(255) default ''::character varying not null,
    agent_id           integer                                    not null
        constraint operation_agent_id_fk
            references agent
            on update cascade on delete cascade,
    total_sales_qty    integer      default 0                     not null,
    total_cost         numeric      default 0.00                  not null,
    total_sales_amount numeric      default 0.00                  not null,
    net_profit         numeric      default 0.00                  not null,
    balance_id         integer                                    not null
        constraint operation_balance_id_fk
            references balance
            on update cascade on delete cascade,
    inventory_id       integer                                    not null
        constraint operation_inventory_id_fk
            references inventory
            on update cascade on delete cascade,
    created_at         timestamp                                  not null,
    updated_at         timestamp                                  not null
);

alter table operation
    owner to root;

create table sale
(
    id           serial
        primary key,
    amount       numeric default 0.0 not null,
    quantity     integer default 0   not null,
    payment_type integer default 1   not null,
    operation_id integer             not null
        constraint sale_operation_id_fk
            references operation
            on update cascade on delete cascade,
    item_id      integer             not null
        constraint sale_item_id_fk
            references item
            on update cascade on delete cascade,
    created_at   timestamp           not null,
    updated_at   timestamp           not null
);

alter table sale
    owner to root;