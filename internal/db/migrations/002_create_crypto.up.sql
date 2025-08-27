CREATE TABLE IF NOT EXISTS crypto (
    symbol TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    current_price DOUBLE PRECISION NOT NULL,
    last_update TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_crypto_last_update ON crypto(last_update);