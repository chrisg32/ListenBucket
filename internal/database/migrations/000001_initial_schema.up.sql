-- Initial schema for ListenBucket
CREATE TABLE IF NOT EXISTS feeds (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    image_url TEXT DEFAULT '',
    is_default INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS episodes (
    id TEXT PRIMARY KEY,
    feed_id TEXT NOT NULL,
    source_id TEXT DEFAULT '',
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    image_url TEXT DEFAULT '',
    audio_url TEXT DEFAULT '',
    duration INTEGER DEFAULT 0,
    source_url TEXT DEFAULT '',
    status TEXT DEFAULT 'pending',
    error_msg TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sources (
    id TEXT PRIMARY KEY,
    feed_id TEXT NOT NULL,
    url TEXT NOT NULL,
    type TEXT DEFAULT 'video',
    title TEXT DEFAULT '',
    image_url TEXT DEFAULT '',
    include_back_catalog INTEGER DEFAULT 1,
    last_checked DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_episodes_feed_id ON episodes(feed_id);
CREATE INDEX IF NOT EXISTS idx_episodes_status ON episodes(status);
CREATE INDEX IF NOT EXISTS idx_sources_feed_id ON sources(feed_id);
