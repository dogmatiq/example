CREATE SCHEMA IF NOT EXISTS bank;

CREATE TABLE IF NOT EXISTS bank.account (
    customer_id UUID NOT NULL,
    account_id  UUID NOT NULL,
    name        TEXT NOT NULL,
    balance     BIGINT NOT NULL DEFAULT 0,

    PRIMARY KEY (customer_id, account_id)
);
