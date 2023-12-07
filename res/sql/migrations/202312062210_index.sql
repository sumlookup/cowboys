-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX IF NOT EXISTS cowboys_id_health_idx ON cowboys(id, health);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX IF EXISTS cowboys_id_health_idx;


