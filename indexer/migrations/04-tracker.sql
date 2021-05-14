
CREATE TABLE IF NOT EXISTS tracks (
    name          text UNIQUE,
    from_addr     text,
    to_addr       text,
    topic         text,
    lastBlockHash text,
    lastBlockNum  numeric,
    startBlock    numeric,
    synced        boolean
);
