-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE urls
ADD user_uuid UUID DEFAULT uuid_generate_v4();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls
DROP COLUMN user_uuid;

DROP EXTENSION "uuid-ossp";
-- +goose StatementEnd
