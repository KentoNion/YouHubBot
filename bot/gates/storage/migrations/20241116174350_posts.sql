-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts(
    id BIGSERIAL PRIMARY KEY,
    source_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    post_source_link TEXT NOT NULL UNIQUE,
    published_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    posted_at TIMESTAMP,
    CONSTRAINT fk_posts_source_id
        FOREIGN KEY (source_id)
            REFERENCES sources (id)
            ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd
