-- +goose Up
-- +goose StatementBegin
CREATE TABLE inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id TEXT NOT NULL,
    class_id TEXT NOT NULL,
    instance_id TEXT NOT NULL,
    app_id INTEGER NOT NULL,
    context_id TEXT NOT NULL,
    amount INTEGER NOT NULL,

    user_id UUID NOT NULL REFERENCES user (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    market_hash_name TEXT NOT NULL,
    market_tradable_restriction INTEGER NOT NULL,
    icon_url TEXT NOT NULL,
    name_color TEXT NOT NULL,

    action_link TEXT NOT NULL,

    tag_type TEXT NOT NULL,
    tag_weapon_internal TEXT NOT NULL,
    tag_weapon_name TEXT NOT NULL,
    tag_quality TEXT NOT NULL,
    tag_rarity TEXT NOT NULL,
    tag_rarity_color TEXT NOT NULL,
    tag_exterior TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
