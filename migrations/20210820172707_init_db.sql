-- +goose Up
CREATE TABLE parcel (
  id BIGSERIAL PRIMARY KEY,
  foo BIGINT NOT NULL
);

-- +goose Down
DROP TABLE parcel;
