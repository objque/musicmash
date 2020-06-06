-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "last_actions" (
    "id"     SERIAL PRIMARY KEY,
    "date"   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "action" VARCHAR(255)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "artists" (
    "id"         SERIAL PRIMARY KEY,
    "name"       VARCHAR(255),
    "poster"     VARCHAR(255),
    "popularity" INTEGER default 0,
    "followers"  INTEGER default 0
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "albums" (
    "id"        SERIAL PRIMARY KEY,
    "artist_id" BIGINT,
    "name"      VARCHAR(255),
    FOREIGN KEY(artist_id) REFERENCES artists(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- +migrate StatementEnd
CREATE UNIQUE INDEX idx_album_art_id_name ON "albums" (artist_id, name);

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "stores" (
    "name" VARCHAR(255) PRIMARY KEY
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "associations" (
    "id"         SERIAL PRIMARY KEY,
    "artist_id"  BIGINT,
    "store_name" VARCHAR(255),
    "store_id"   VARCHAR(255),
    FOREIGN KEY(artist_id)  REFERENCES artists(id)  ON UPDATE RESTRICT ON DELETE RESTRICT,
    FOREIGN KEY(store_name) REFERENCES stores(name) ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- +migrate StatementEnd
CREATE UNIQUE INDEX idx_art_store_name_id ON "associations" (artist_id, store_name, store_id);

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "releases" (
    "id"         SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "artist_id"  BIGINT,
    "title"      VARCHAR(1000),
    "poster"     VARCHAR(255),
    "released"   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "store_name" VARCHAR(255),
    "store_id"   VARCHAR(255),
    FOREIGN KEY(artist_id)  REFERENCES artists(id)  ON UPDATE RESTRICT ON DELETE RESTRICT,
    FOREIGN KEY(store_name) REFERENCES stores(name) ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- +migrate StatementEnd
CREATE UNIQUE INDEX idx_rel_store_name_store_id ON "releases" (store_name, store_id);
CREATE INDEX idx_releases_created ON "releases" (created_at);

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "subscriptions" (
    "id"         SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "user_name"  VARCHAR(255),
    "artist_id"  BIGINT,
    FOREIGN KEY(artist_id) REFERENCES artists(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- +migrate StatementEnd
CREATE UNIQUE INDEX idx_user_name_artist_id ON "subscriptions" (user_name, artist_id);

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "notifications" (
    "id"         SERIAL PRIMARY KEY,
    "date"       TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "user_name"  VARCHAR(255),
    "release_id" BIGINT,
    "is_coming"  bool DEFAULT false,
    FOREIGN KEY(release_id) REFERENCES releases(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- +migrate StatementEnd
CREATE UNIQUE INDEX idx_user_name_release_id_is_coming ON "notifications" (user_name, release_id, is_coming);

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "notification_services" (
    "id" VARCHAR(255) PRIMARY KEY
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "notification_settings" (
    "id"        SERIAL PRIMARY KEY,
    "user_name" VARCHAR(255),
    "service"   VARCHAR(255),
    "data"      VARCHAR(255),
    FOREIGN KEY(service) REFERENCES notification_services(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- +migrate StatementEnd
CREATE UNIQUE INDEX idx_user_name_service ON "notification_settings" (user_name, service);

-- +migrate Down
DROP TABLE IF EXISTS "last_actions";
DROP TABLE IF EXISTS "notifications";
DROP TABLE IF EXISTS "notification_settings";
DROP TABLE IF EXISTS "notification_services";
DROP TABLE IF EXISTS "releases";
DROP TABLE IF EXISTS "albums";
DROP TABLE IF EXISTS "subscriptions";
DROP TABLE IF EXISTS "associations";
DROP TABLE IF EXISTS "stores";
DROP TABLE IF EXISTS "artists";
