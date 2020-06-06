-- +migrate Up
CREATE INDEX idx_releases_artists_id ON "releases" (artist_id);

-- +migrate Down
DROP INDEX idx_releases_artists_id;
