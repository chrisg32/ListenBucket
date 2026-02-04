-- Rollback initial schema
DROP INDEX IF EXISTS idx_sources_feed_id;
DROP INDEX IF EXISTS idx_episodes_status;
DROP INDEX IF EXISTS idx_episodes_feed_id;
DROP TABLE IF EXISTS sources;
DROP TABLE IF EXISTS episodes;
DROP TABLE IF EXISTS feeds;
