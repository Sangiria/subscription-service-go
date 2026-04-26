-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_name TEXT NOT NULL,
    cost INTEGER NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);
CREATE INDEX idx_subscriptions_user_period ON subscriptions (
    user_id,
    subscription_name,
    start_date
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_subscriptions_user_period;
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd