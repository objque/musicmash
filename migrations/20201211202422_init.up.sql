BEGIN;

CREATE TABLE IF NOT EXISTS "artists" (
    "id"          serial       NOT NULL,
    "created_at"  timestamp    NOT NULL DEFAULT now(),
    "name"        varchar(255) NOT NULL,
    "poster"      varchar(255),
    "is_verified" boolean      NOT NULL DEFAULT FALSE,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX "idx_artists_name_poster" ON "artists" USING BTREE ("name","poster");

CREATE TABLE IF NOT EXISTS "artist_associations" (
    "id"         serial       NOT NULL,
    "artist_id"  int          NOT NULL,
    "store_name" varchar(255) NOT NULL,
    "store_id"   varchar(255) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (artist_id) REFERENCES "artists"(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE INDEX "idx_artist_associations_store_id" ON "artist_associations" USING BTREE ("store_id");
CREATE UNIQUE INDEX "idx_artist_associations_store_name_store_id" ON "artist_associations" USING BTREE ("artist_id", "store_name","store_id");

CREATE TABLE IF NOT EXISTS "artist_details" (
    "id"         serial NOT NULL,
    "artist_id"  int    NOT NULL,
    "genres"     text[] NOT NULL DEFAULT '{}',
    "bio"        text,
    "popularity" int4,
    "followers"  int4,
    "listeners"  int4,
    "world_rank" int4,
    PRIMARY KEY (id),
    FOREIGN KEY (artist_id) REFERENCES "artists"(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE UNIQUE INDEX "idx_artist_details_artist_id" ON "artist_details" USING BTREE ("artist_id");

CREATE TABLE IF NOT EXISTS "artist_headers" (
    "id"         serial    NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "artist_id"  int       NOT NULL,
    "width"      int       NOT NULL DEFAULT 0,
    "height"     int       NOT NULL DEFAULT 0,
    "url"        text      NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (artist_id) REFERENCES "artists"(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE UNIQUE INDEX "idx_artist_headers_artist_id_url" ON "artist_headers" USING BTREE ("artist_id","url");

CREATE TABLE IF NOT EXISTS "artist_gallery" (
    "id"         serial    NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "artist_id"  int       NOT NULL,
    "width"      int       NOT NULL DEFAULT 0,
    "height"     int       NOT NULL DEFAULT 0,
    "url"        text      NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (artist_id) REFERENCES "artists"(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE UNIQUE INDEX "idx_artist_gallery_artist_id_url" ON "artist_gallery" USING BTREE ("artist_id","url");

CREATE TABLE IF NOT EXISTS "artist_posters" (
    "id"         serial    NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "artist_id"  int       NOT NULL,
    "width"      int       NOT NULL DEFAULT 0,
    "height"     int       NOT NULL DEFAULT 0,
    "url"        text      NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (artist_id) REFERENCES "artists"(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE UNIQUE INDEX "idx_artist_posters_artist_id_url" ON "artist_posters" USING BTREE ("artist_id","url");

CREATE TABLE IF NOT EXISTS "artist_external_links" (
    "id"         serial    NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "artist_id"  int       NOT NULL,
    "service"    text      NOT NULL,
    "url"        text      NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (artist_id) REFERENCES "artists"(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE UNIQUE INDEX "idx_artist_external_links_artist_id_service_url" ON "artist_external_links" USING BTREE ("artist_id","service","url");

CREATE TABLE IF NOT EXISTS "subscriptions" (
    "id"         serial       NOT NULL,
    "created_at" timestamp    NOT NULL DEFAULT now(),
    "user_name"  varchar(255) NOT NULL,
    "artist_id"  bigint       NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY(artist_id) REFERENCES artists(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE UNIQUE INDEX idx_user_name_artist_id ON "subscriptions" (user_name, artist_id);

CREATE TABLE IF NOT EXISTS "releases" (
    "id"           serial        NOT NULL,
    "created_at"   timestamp     NOT NULL DEFAULT now(),
    "artist_id"    int           NOT NULL,
    "type"         varchar(32)   NOT NULL,
    "tracks_count" integer       NOT NULL,
    "duration_ms"  bigint        NOT NULL,
    "title"        varchar(1000) NOT NULL,
    "released"     timestamp     NOT NULL,
    "is_explicit"  boolean       NOT NULL DEFAULT false,
    "popularity"   integer       DEFAULT NULL,
    "poster"       varchar(255),
    "spotify_id"   varchar(32)   NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (artist_id) REFERENCES "artists"(id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE INDEX "idx_artist_releases_artist_id" ON "releases" USING BTREE ("artist_id");
CREATE INDEX "idx_artist_releases_created_at" ON "releases" USING BTREE ("created_at");
CREATE UNIQUE INDEX "idx_artist_releases_artist_id_id" ON "releases" USING BTREE ("artist_id", "id");

CREATE TABLE IF NOT EXISTS "stores" (
    "name" varchar(255),
    PRIMARY KEY (name)
);
INSERT INTO "stores" (name) VALUES ('spotify') ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS "last_actions" (
    "id"     serial    NOT NULL,
    "date"   timestamp NOT NULL DEFAULT now(),
    "action" varchar(255),
    PRIMARY KEY (id)
);

COMMIT;