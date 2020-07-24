-- +migrate Up
INSERT INTO "stores" (name) VALUES ('itunes') ON CONFLICT DO NOTHING;

-- +migrate Down
