-- +goose Up
-- +goose StatementBegin
ALTER TABLE urls
ADD is_deleted BOOL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls
DROP COLUMN is_deleted;
-- +goose StatementEnd
