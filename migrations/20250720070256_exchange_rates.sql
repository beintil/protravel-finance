-- +goose Up
-- +goose StatementBegin
CREATE TABLE exchange_rates (
    id UUID PRIMARY KEY,
    base_currency VARCHAR(3) NOT NULL REFERENCES currencies(code),
    target_currency VARCHAR(3) NOT NULL REFERENCES currencies(code),
    rate BIGINT NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (base_currency, target_currency, date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
