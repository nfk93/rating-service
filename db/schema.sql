CREATE TABLE users (
  id   text NOT NULL PRIMARY KEY,
  name text NOT NULL
);

CREATE TABLE games (
  id text NOT NULL PRIMARY KEY,
  name text NOT NULL
);

CREATE TABLE ratings (
  user_id text references users(id),
  game_id text references games(id),
  rating int NOT NULL,
  PRIMARY KEY(user_id, game_id)
);