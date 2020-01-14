-- +migrate Up
ALTER TABLE "albums" RENAME TO "_albums";

-- +migrate Down
ALTER TABLE "_albums" RENAME TO "albums";
