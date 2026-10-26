CREATE TABLE IF NOT EXISTS accounts (
    id      VARCHAR(64)  PRIMARY KEY,
    email   VARCHAR(255) NOT NULL,
    balance BIGINT       NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions (
    id         VARCHAR(64) PRIMARY KEY,
    account_id VARCHAR(64) NOT NULL,
    amount     BIGINT      NOT NULL,
    type       VARCHAR(32) NOT NULL,
    INDEX idx_transactions_account_id (account_id)
);