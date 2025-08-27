-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    steam_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    cash NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    trade_url TEXT NOT NULL,
    avatar_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_users_steam_id ON users (steam_id);
CREATE INDEX idx_users_email ON users (email);

CREATE TYPE offer_status AS ENUM ('onsale', 'reserved', 'sold', 'canceled');
CREATE TABLE offers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id TEXT REFERENCES users (steam_id),
    bot_steam_id TEXT DEFAULT 0,
    price DOUBLE PRECISION NOT NULL,
    steam_trade_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    status offer_status NOT NULL DEFAULT 'onsale',
    reserved_until TIMESTAMPTZ,

    asset_id TEXT NOT NULL,
    class_id TEXT NOT NULL,
    instance_id TEXT NOT NULL,
    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    market_tradable_restriction INTEGER NOT NULL,
    icon_url TEXT NOT NULL,
    name_color TEXT NOT NULL,
    action_link TEXT,

    tag_type TEXT NOT NULL,
    tag_weapon_internal TEXT NOT NULL,
    tag_weapon_name TEXT NOT NULL,
    tag_quality TEXT NOT NULL,
    tag_rarity TEXT NOT NULL,
    tag_rarity_color TEXT NOT NULL,
    tag_exterior TEXT NOT NULL
);

CREATE INDEX idx_offers_seller_id ON offers (seller_id);
CREATE INDEX idx_offers_asset_id ON offers (asset_id);
CREATE INDEX idx_offers_status ON offers (status);
CREATE INDEX idx_steam_trade_id ON offers (steam_trade_id);
-- CREATE INDEX idx_offers_reserved_until ON offers (reserved_until);


CREATE TYPE transaction_status AS ENUM ('completed', 'failed');
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    offer_id UUID NOT NULL REFERENCES offers (id),
    seller_id TEXT NOT NULL,
    buyer_id TEXT NOT NULL,
    status transaction_status NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    market_tradable_restriction INTEGER NOT NULL,
    icon_url TEXT NOT NULL,
    name_color TEXT NOT NULL,
    action_link TEXT,

    tag_type TEXT NOT NULL,
    tag_weapon_internal TEXT NOT NULL,
    tag_weapon_name TEXT NOT NULL,
    tag_quality TEXT NOT NULL,
    tag_rarity TEXT NOT NULL,
    tag_rarity_color TEXT NOT NULL,
    tag_exterior TEXT NOT NULL
);

CREATE INDEX idx_transactions_offer_id ON transactions (offer_id);
CREATE INDEX idx_transactions_seller_id ON transactions (seller_id);
CREATE INDEX idx_transactions_buyer_id ON transactions (buyer_id);
CREATE INDEX idx_transactions_status ON transactions (status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS offers;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS offer_status;
DROP TYPE IF EXISTS transaction_status;
-- +goose StatementEnd
