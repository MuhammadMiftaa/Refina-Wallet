package data

import "time"

var (
	DEVELOPMENT_MODE = "development"
	STAGING_MODE     = "staging"
	PRODUCTION_MODE  = "production"

	OUTBOX_PUBLISH_EXCHANGE     = "refina_microservice"
	OUTBOX_PUBLISH_INTERVAL     = 5 * time.Second
	OUTBOX_PUBLISH_BATCH        = 100
	OUTBOX_PUBLISH_MAX_RETRIES  = 5
	OUTBOX_EVENT_WALLET_CREATED = "wallet.created"
	OUTBOX_EVENT_WALLET_UPDATED = "wallet.updated"
	OUTBOX_EVENT_WALLET_DELETED = "wallet.deleted"

	INITIAL_DEPOSIT_CATEGORY_ID = "00000000-0000-0000-0000-000000000000"
	INITIAL_DEPOSIT_DESC        = "Deposit awal"
)
