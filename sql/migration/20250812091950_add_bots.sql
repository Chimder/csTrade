-- +goose Up
-- +goose StatementBegin
CREATE TABLE bots (
    stream_id BIGINT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    shared_secret TEXT NOT NULL,
    skin_count INTEGER NOT NULL DEFAULT 0,
    identity_secret TEXT NOT NULL,
    device_id TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bots;
-- +goose StatementEnd
