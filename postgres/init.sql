create table currency_rates
(
    id              serial            not null        primary key,
    code            varchar(6)      not null,
    value           numeric(20, 10) not null,
    last_updated_at timestamp,
    batch_id        uuid            not null,
    inserted_at     timestamp default now(),
    constraint currency_rates_pk unique (code, last_updated_at)
);

create table api_calls
(
    id            uuid     not null        primary key,
    status_code   integer  not null,
    response_time interval not null,
    timeout       boolean   default false,
    error_message text,
    called_at     timestamp default now()
);