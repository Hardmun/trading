CREATE TABLE IF NOT EXISTS users
(
    "id"
        INTEGER PRIMARY KEY
        AUTOINCREMENT
        NOT NULL,
    "user"
        TEXT
        NOT NULL,
    "password"
        TEXT
        NOT NULL,
    "active"
        BOOLEAN,
    "requestID"
        TEXT
        NOT NULL
);

PRAGMA journal_mode= WAL;