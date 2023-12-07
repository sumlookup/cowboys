-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

INSERT INTO cowboys (name, health, damage) VALUES
    ('John', 10, 1),
    ('Bill', 8, 2),
    ('Sam', 10, 1),
    ('Peter', 5, 3),
    ('Philip', 15, 1) ON CONFLICT DO NOTHING;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back