CREATE TABLE IF NOT EXISTS suggestions
(
    id  INTEGER PRIMARY KEY,
    txt TEXT NOT NULL UNIQUE,
    txt_answer TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_id ON suggestions (id);

CREATE TABLE IF NOT EXISTS complaints
(
    id  INTEGER PRIMARY KEY,
    txt TEXT NOT NULL UNIQUE,
    txt_answer TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_id ON complaints (id);

CREATE TABLE IF NOT EXISTS apps
(
    id     INTEGER PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);