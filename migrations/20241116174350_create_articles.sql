-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE posts(
    post_ID INT NOT NULL,
    tg_chan_name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    post_source_link VARCHAR(255) NOT NULL,
    published_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    posted_at TIMESTAMP,
    PRIMARY KEY(tg_chan_name, post_ID),
    FOREIGN KEY (tg_chan_name) REFERENCES tg_sources(tg_chan_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE posts;
-- +goose StatementEnd
