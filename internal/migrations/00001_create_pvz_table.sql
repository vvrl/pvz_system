-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pvz (
    id SERIAL PRIMARY KEY,
    city TEXT NOT NULL,
    registrationDate TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pvz;
-- +goose StatementEnd
