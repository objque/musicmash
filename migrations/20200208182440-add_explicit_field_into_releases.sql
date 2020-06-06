-- +migrate Up
ALTER TABLE "releases" ADD COLUMN "explicit" BOOLEAN NOT NULL DEFAULT false;

-- +migrate Down
ALTER TABLE "releases" DROP COLUMN "explicit";
