-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE wallet_types ADD COLUMN description TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE wallet_types DROP COLUMN description;
-- +goose StatementEnd
