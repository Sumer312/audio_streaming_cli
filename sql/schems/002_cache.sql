-- +goose Up 
CREATE TABLE Cache (
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  song_id INTEGER NOT NULL,
  cached_url TEXT UNIQUE,
  FOREIGN KEY(songs_id) REFERENCES Songs(id)
);

-- +goose Down
DROP TABLE Cache;
