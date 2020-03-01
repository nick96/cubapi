-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE users (
         id SERIAL              PRIMARY KEY
       , email     VARCHAR(256) NOT NULL UNIQUE
       , firstName VARCHAR(256) NOT NULL
       , lastName  VARCHAR(256) NOT NULL
       , password  VARCHAR(100) NOT NULL
       , salt      VARCHAR(15)  NOT NULL UNIQUE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;
