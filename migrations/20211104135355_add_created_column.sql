-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE parcel ADD COLUMN created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE parcel DROP COLUMN created;
