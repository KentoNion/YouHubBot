-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE admins(
    admin varchar(255) NOT NULL,
    role varchar(255) NoT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE admins;
-- +goose StatementEnd
