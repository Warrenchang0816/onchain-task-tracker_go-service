-- Optional migration bookkeeping
CREATE TABLE IF NOT EXISTS schema_migrations (
    version      VARCHAR(100) PRIMARY KEY,
    applied_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Core tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    task_id TEXT UNIQUE,
    wallet_address TEXT,
    assignee_wallet_address TEXT,
    title TEXT NOT NULL DEFAULT 'Untitled task',
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
    priority VARCHAR(50) NOT NULL DEFAULT 'MEDIUM',
    reward_amount NUMERIC,
    fee_bps INTEGER,
    onchain_status VARCHAR(50),
    fund_tx_hash VARCHAR(255),
    due_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS task_id TEXT;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS wallet_address TEXT;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS assignee_wallet_address TEXT;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS title TEXT;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'OPEN';
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS priority VARCHAR(50) NOT NULL DEFAULT 'MEDIUM';
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS reward_amount NUMERIC;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS fee_bps INTEGER;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS onchain_status VARCHAR(50);
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS fund_tx_hash VARCHAR(255);
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS due_date TIMESTAMPTZ;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

UPDATE tasks
SET title = COALESCE(NULLIF(title, ''), 'Untitled task')
WHERE title IS NULL OR title = '';

ALTER TABLE tasks ALTER COLUMN title SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);

-- Blockchain logs table
CREATE TABLE IF NOT EXISTS task_blockchain_logs (
    id BIGSERIAL PRIMARY KEY,
    task_id TEXT NOT NULL,
    action VARCHAR(100) NOT NULL,
    tx_hash VARCHAR(255),
    chain_id VARCHAR(50),
    contract_address VARCHAR(255),
    status VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE task_blockchain_logs ADD COLUMN IF NOT EXISTS task_id TEXT;
ALTER TABLE task_blockchain_logs ADD COLUMN IF NOT EXISTS action VARCHAR(100);
ALTER TABLE task_blockchain_logs ADD COLUMN IF NOT EXISTS tx_hash VARCHAR(255);
ALTER TABLE task_blockchain_logs ADD COLUMN IF NOT EXISTS chain_id VARCHAR(50);
ALTER TABLE task_blockchain_logs ADD COLUMN IF NOT EXISTS contract_address VARCHAR(255);
ALTER TABLE task_blockchain_logs ADD COLUMN IF NOT EXISTS status VARCHAR(50);
ALTER TABLE task_blockchain_logs ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_task_blockchain_logs_task_id ON task_blockchain_logs(task_id);
CREATE INDEX IF NOT EXISTS idx_task_blockchain_logs_action ON task_blockchain_logs(action);
CREATE INDEX IF NOT EXISTS idx_task_blockchain_logs_created_at ON task_blockchain_logs(created_at DESC);

-- NFT orders table
CREATE TABLE IF NOT EXISTS nft_orders (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    image TEXT,
    price NUMERIC,
    recipient_wallet TEXT,
    creator_wallet TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_nft_orders_creator_wallet ON nft_orders(creator_wallet);
CREATE INDEX IF NOT EXISTS idx_nft_orders_recipient_wallet ON nft_orders(recipient_wallet);
CREATE INDEX IF NOT EXISTS idx_nft_orders_created_at ON nft_orders(created_at DESC);

INSERT INTO schema_migrations(version)
VALUES ('001_init_schema')
ON CONFLICT (version) DO NOTHING;