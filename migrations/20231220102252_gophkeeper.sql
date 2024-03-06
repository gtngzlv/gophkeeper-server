-- +goose Up
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL NOT NULL PRIMARY KEY,
                                     email TEXT NOT NULL UNIQUE,
                                     password_hash BYTEA NOT NULL,
                                     secret_key_hash BYTEA NOT NULL,
                                     encrypted_key BYTEA NOT NULL,
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_email on users(email);

CREATE TABLE IF NOT EXISTS PERSONAL_DATA(
    ID SERIAL PRIMARY KEY,
    PDATA TEXT NOT NULL,
    USER_ID INT NOT NULL,
    FOREIGN KEY(USER_ID) REFERENCES USERS(ID),
    CREATED_AT TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL);

-- +goose Down
-- +goose StatementBegin
DROP TABLE PERSONAL_DATA;
DROP TABLE USERS;
-- +goose StatementEnd
