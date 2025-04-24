-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
  uuid SERIAL PRIMARY KEY,
  short_url VARCHAR(255) NOT NULL,
  original_url VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS original_url ON urls (original_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls;
-- +goose StatementEnd