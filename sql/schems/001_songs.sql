-- +goose Up 

CREATE TABLE songs (
  title VARCHAR(255) PRIMARY KEY,
  url VARCHAR(1000) UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
);

-- +goose Down
DROP TABLE songs;
