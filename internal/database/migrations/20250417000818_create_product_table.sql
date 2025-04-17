-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    dateTime TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    type TEXT NOT NULL,
    receptionId TEXT REFERENCES receptions(id) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
