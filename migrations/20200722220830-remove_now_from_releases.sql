-- +migrate Up
ALTER TABLE "releases" ALTER COLUMN "released" DROP DEFAULT;

-- +migrate Down
ALTER TABLE "releases" ALTER COLUMN "released" SET DEFAULT CURRENT_TIMESTAMP;

