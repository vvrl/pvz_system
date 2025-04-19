-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS receptions (
    id SERIAL PRIMARY KEY,
    dateTime TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    pvzId INTEGER REFERENCES pvz(id) NOT NULL,
    status TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE receptions;
-- +goose StatementEnd
