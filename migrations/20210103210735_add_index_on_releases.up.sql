BEGIN;

CREATE INDEX IF NOT EXISTS "idx_releases_released" on "releases" USING BTREE ("released");

COMMIT;