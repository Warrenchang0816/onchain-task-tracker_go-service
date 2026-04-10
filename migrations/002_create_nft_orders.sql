CREATE TABLE IF NOT EXISTS nft_orders (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    image TEXT DEFAULT '',
    price TEXT DEFAULT '',
    creator_wallet TEXT DEFAULT '',
    recipient_wallet TEXT DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);