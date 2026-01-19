-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW view_user_wallets AS
SELECT 
    wallets.id, users.id AS user_id,
	wallets.number AS wallet_number, wallets.balance AS wallet_balance,
	wallets.name AS wallet_name, wallet_types.name AS wallet_type_name,
	wallet_types.type AS wallet_type
FROM wallets
LEFT JOIN users ON users.id = wallets.user_id AND users.deleted_at IS NULL
LEFT JOIN wallet_types ON wallet_types.id = wallets.wallet_type_id AND wallet_types.deleted_at IS NULL
WHERE wallets.deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS view_user_wallets;
-- +goose StatementEnd
