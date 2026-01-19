-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW view_user_wallets_group_by_type AS
SELECT 
	wallets.user_id,
	wallet_types.type AS type,
	JSON_AGG(
		JSON_BUILD_OBJECT(
			'id', wallets.id,
			'name', wallets.name,
			'number', wallets.number,
			'balance', wallets.balance
		)
	) AS wallets 
FROM wallets
JOIN wallet_types ON wallet_types.id = wallets.wallet_type_id AND wallet_types.deleted_at IS NULL
WHERE wallets.deleted_at IS NULL
GROUP BY wallets.user_id, wallet_types.type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS view_user_wallets_group_by_type;
-- +goose StatementEnd
