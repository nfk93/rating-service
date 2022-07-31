CREATE TYPE rating_system_enum AS ENUM ('glicko', 'elo');

CREATE TABLE users (
  id   uuid NOT NULL PRIMARY KEY,
  name string NOT NULL
);

CREATE TABLE games (
  id uuid NOT NULL PRIMARY KEY,
  name string NOT NULL,
  rating_system rating_system_enum
);

CREATE TABLE elo_rating (
  user_id uuid references users(id),
  game_id uuid references users(id),

  rating int NOT NULL,

  PRIMARY KEY(user_id, game_id)
);

CREATE TABLE matches (
  id uuid NOT NULL PRIMARY KEY,
  game_id uuid references games(id),
  happened_at timestamp with time zone NOT NULL
);

CREATE TABLE match_player (
  match_id uuid references matches(id),
  user_id uuid references users(id),
  is_winner boolean NOT NULL,
  score integer,
  PRIMARY KEY(match_id, user_id)
);

