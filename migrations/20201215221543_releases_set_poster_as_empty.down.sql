BEGIN;

ALTER TABLE releases ALTER COLUMN poster DROP NOT NULL;

UPDATE releases SET poster = NULL WHERE poster = '';

COMMIT;