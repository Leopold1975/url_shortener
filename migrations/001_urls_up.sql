-- +goose up
CREATE TABLE IF NOT EXISTS urls (
uuid text primary key,
url_long text not null unique,
url_short text not null unique,
created_at timestamptz not null DEFAULT current_timestamp,
clicks integer);

-- +goose down
DROP TABLE IF EXISTS urls;