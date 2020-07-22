-- +migrate Up
ALTER TABLE "artists"       ALTER COLUMN "name"       SET NOT NULL;

ALTER TABLE "associations"  ALTER COLUMN "artist_id"  SET NOT NULL;
ALTER TABLE "associations"  ALTER COLUMN "store_name" SET NOT NULL;
ALTER TABLE "associations"  ALTER COLUMN "store_id"   SET NOT NULL;

ALTER TABLE "last_actions"  ALTER COLUMN "action"     SET NOT NULL;

ALTER TABLE "releases"      ALTER COLUMN "artist_id"  SET NOT NULL;
ALTER TABLE "releases"      ALTER COLUMN "title"      SET NOT NULL;
ALTER TABLE "releases"      ALTER COLUMN "released"   SET NOT NULL;
ALTER TABLE "releases"      ALTER COLUMN "store_name" SET NOT NULL;
ALTER TABLE "releases"      ALTER COLUMN "store_id"   SET NOT NULL;

ALTER TABLE "subscriptions" ALTER COLUMN "user_name"  SET NOT NULL;
ALTER TABLE "subscriptions" ALTER COLUMN "artist_id"  SET NOT NULL;

-- +migrate Down
ALTER TABLE "artists"       ALTER COLUMN "name"       DROP NOT NULL;

ALTER TABLE "associations"  ALTER COLUMN "artist_id"  DROP NOT NULL;
ALTER TABLE "associations"  ALTER COLUMN "store_name" DROP NOT NULL;
ALTER TABLE "associations"  ALTER COLUMN "store_id"   DROP NOT NULL;

ALTER TABLE "last_actions"  ALTER COLUMN "action"     DROP NOT NULL;

ALTER TABLE "releases"      ALTER COLUMN "artist_id"  DROP NOT NULL;
ALTER TABLE "releases"      ALTER COLUMN "title"      DROP NOT NULL;
ALTER TABLE "releases"      ALTER COLUMN "released"   DROP NOT NULL;
ALTER TABLE "releases"      ALTER COLUMN "store_name" DROP NOT NULL;
ALTER TABLE "releases"      ALTER COLUMN "store_id"   DROP NOT NULL;

ALTER TABLE "subscriptions" ALTER COLUMN "user_name"  DROP NOT NULL;
ALTER TABLE "subscriptions" ALTER COLUMN "artist_id"  DROP NOT NULL;
