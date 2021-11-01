-- +goose Up
CREATE TABLE workplace (
  id BIGSERIAL PRIMARY KEY,
  foo BIGINT NOT NULL
);

-- +goose Down
DROP TABLE workplace;
