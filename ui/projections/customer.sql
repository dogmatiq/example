-- customers contains one row per bank customer.
--
-- It is populated by the "customers" projection, implemented by the
-- CustomerProjectionHandler type in customer.go.
CREATE TABLE IF NOT EXISTS customers (
    id   TEXT NOT NULL, -- unique customer identifier
    name TEXT NOT NULL, -- customer's full name

    PRIMARY KEY (id)
);
