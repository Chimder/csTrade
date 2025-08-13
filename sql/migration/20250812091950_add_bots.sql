-- +goose Up
-- +goose StatementBegin
CREATE TABLE bots (
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    stream_id TEXT NOT NULL,
    shared_secret TEXT NOT NULL,
    identity_secret TEXT NOT NULL,
    device_id TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bots;
-- +goose StatementEnd
