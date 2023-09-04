CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    username text NOT NULL UNIQUE,
    password text NOT NULL
);