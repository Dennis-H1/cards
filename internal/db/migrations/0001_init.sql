CREATE TABLE cards (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    front      TEXT NOT NULL,
    back       TEXT NOT NULL,
    source     TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tags (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    name                TEXT NOT NULL UNIQUE,
    overview            TEXT,
    overview_updated_at TIMESTAMP
);

CREATE TABLE card_tags (
    card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    tag_id  INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (card_id, tag_id)
);

CREATE INDEX idx_card_tags_tag_id ON card_tags(tag_id);

CREATE TABLE reviews (
    card_id          INTEGER PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    ease_factor      REAL NOT NULL DEFAULT 2.5,
    interval_days    INTEGER NOT NULL DEFAULT 0,
    repetitions      INTEGER NOT NULL DEFAULT 0,
    due_at           TIMESTAMP NOT NULL,
    review_count     INTEGER NOT NULL DEFAULT 0,
    last_reviewed_at TIMESTAMP
);

CREATE INDEX idx_reviews_due_at ON reviews(due_at);

CREATE TABLE events (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    type       TEXT NOT NULL CHECK (type IN ('card_created', 'card_edited')),
    card_id    INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    seen_at    TIMESTAMP
);

CREATE INDEX idx_events_seen_at ON events(seen_at);
