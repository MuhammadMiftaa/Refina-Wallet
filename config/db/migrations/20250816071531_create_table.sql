-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallet_types (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    deleted_at timestamptz,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL
);

CREATE TABLE wallets (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    deleted_at timestamptz,
    user_id uuid NOT NULL,
    wallet_type_id uuid NOT NULL,
    name VARCHAR(50) NOT NULL,
    balance numeric(18,2) NOT NULL,
    number VARCHAR(50) NOT NULL,
    FOREIGN KEY (wallet_type_id) REFERENCES wallet_types(id) ON DELETE RESTRICT
);

-- Create indexes for WHERE operations
CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_wallets_wallet_type_id ON wallets(wallet_type_id);
CREATE INDEX idx_wallets_deleted_at ON wallets(deleted_at);
CREATE INDEX idx_wallet_types_deleted_at ON wallet_types(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop indexes first
DROP INDEX IF EXISTS idx_wallet_types_deleted_at;
DROP INDEX IF EXISTS idx_wallets_deleted_at;
DROP INDEX IF EXISTS idx_wallets_wallet_type_id;
DROP INDEX IF EXISTS idx_wallets_user_id;

-- Drop tables
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS wallet_types;
-- +goose StatementEnd
