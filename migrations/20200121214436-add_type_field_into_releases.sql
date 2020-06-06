-- +migrate Up
ALTER TABLE "releases" ADD COLUMN "type" VARCHAR(20) NOT NULL DEFAULT 'album';

-- +migrate Down
ALTER TABLE "releases" DROP COLUMN "type";
