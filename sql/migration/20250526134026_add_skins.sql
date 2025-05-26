-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE skins (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    category_id TEXT,
    category_name TEXT,
    pattern_id TEXT,
    pattern_name TEXT,
    min_float DOUBLE PRECISION,
    max_float DOUBLE PRECISION,
    rarity_id TEXT,
    rarity_name TEXT,
    rarity_color TEXT,
    stattrak BOOLEAN,
    souvenir BOOLEAN,
    paint_index TEXT,
    legacy_model BOOLEAN,
    image TEXT NOT NULL,
    phase TEXT,
    team_id TEXT,
    team_name TEXT
);

CREATE INDEX ON skins (name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS skins
-- +goose StatementEnd
