-- +goose Up 
CREATE TABLE Cache (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  song_id INTEGER NOT NULL,
  cached_url TEXT UNIQUE,
  FOREIGN KEY(song_id) REFERENCES Songs(id)
);

-- +goose Down
DROP TABLE Cache;
