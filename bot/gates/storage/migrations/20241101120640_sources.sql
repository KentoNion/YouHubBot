-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sources(
    id SERIAL PRIMARY KEY,
    tg_chan_name VARCHAR(255) NOT NULL,
    tg_chan_link VARCHAR(255) NOT NULL,
    --source_chan VARCHAR(255) NOT NULL,
    source_chan_link VARCHAR(255) NOT NULL UNIQUE,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sources;
-- +goose StatementEnd
