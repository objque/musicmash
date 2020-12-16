BEGIN;

UPDATE releases SET poster = '' WHERE poster IS NULL;

ALTER TABLE releases ALTER COLUMN poster SET NOT NULL;

COMMIT;