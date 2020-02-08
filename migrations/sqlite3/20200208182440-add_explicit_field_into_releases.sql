-- +migrate Up
ALTER TABLE "releases" ADD "explicit" BOOLEAN NOT NULL DEFAULT 0;

-- +migrate Down
PRAGMA foreign_keys=off;

-- +migrate StatementBegin
CREATE TABLE "temp_releases" (
    "id"         INTEGER PRIMARY KEY AUTOINCREMENT,
    "created_at" DATETIME,
    "artist_id"  BIGINT,
    "title"      VARCHAR(1000),
    "poster"     VARCHAR(255),
    "released"   DATETIME,
    "store_name" VARCHAR(255),
    "store_id"   VARCHAR(255),
    "type"       VARCHAR(20),
    FOREIGN KEY(artist_id)  REFERENCES artists(id)  ON UPDATE RESTRICT ON DELETE RESTRICT,
    FOREIGN KEY(store_name) REFERENCES stores(name) ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- +migrate StatementEnd

-- +migrate StatementBegin
INSERT INTO "temp_releases" (id,created_at,artist_id,title,poster,released,store_name,store_id,type)
SELECT id,created_at,artist_id,title,poster,released,store_name,store_id,type FROM releases;
-- +migrate StatementEnd

DROP TABLE "releases";
ALTER TABLE "temp_releases" RENAME TO "releases";

CREATE UNIQUE INDEX idx_rel_store_name_store_id ON "releases" (store_name, store_id);
CREATE INDEX idx_releases_created ON "releases" (created_at);
CREATE INDEX idx_releases_artists_id ON "releases" (artist_id);

PRAGMA foreign_keys=on;