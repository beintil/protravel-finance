-- +goose Up
-- +goose StatementBegin
CREATE TABLE currencies (
    code VARCHAR(3) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);
INSERT INTO currencies (code, name) VALUES ('USD', 'US Dollar');
INSERT INTO currencies (code, name) VALUES ('EUR', 'Euro');
INSERT INTO currencies (code, name) VALUES ('RUB', 'Russian Ruble');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
