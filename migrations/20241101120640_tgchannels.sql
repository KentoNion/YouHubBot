-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS tg_sources(
    tg_chan_name VARCHAR(255) NOT NULL,
    tg_chan_link VARCHAR(255) NOT NULL,
    source_chan VARCHAR(255) NOT NULL,
    source_chan_link VARCHAR(255) NOT NULL,
    primary key(tg_chan_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE tg_sources;
-- +goose StatementEnd
