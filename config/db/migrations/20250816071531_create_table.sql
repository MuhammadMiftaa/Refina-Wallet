-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallet_types (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL
);

CREATE TABLE wallets (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id uuid NOT NULL,
    wallet_type_id uuid NOT NULL,
    name VARCHAR(50) NOT NULL,
    balance numeric(18,2) NOT NULL,
    number VARCHAR(50) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS wallet_types;
-- +goose StatementEnd
