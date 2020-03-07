-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE attendance (
  id              SERIAL       PRIMARY KEY
  , cubName       VARCHAR(256) NOT NULL
  , parentSignIn  TEXT         NOT NULL
  , parentSignOut TEXT         NOT NULL
  , approvingUser  VARCHAR(256) NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE attendance;
