CREATE TYPE rating_system_enum AS ENUM ('glicko', 'elo');

CREATE TABLE users (
  id   uuid NOT NULL PRIMARY KEY,
  name text NOT NULL
);

CREATE TABLE games (
  id uuid NOT NULL PRIMARY KEY,
  name text NOT NULL,
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
  game_id uuid NOT NULL references games(id),
  ratings_updated boolean NOT NULL default false,
  is_finished boolean NOT NULL default false,
  happened_at timestamp with time zone NOT NULL
);

CREATE TABLE match_player (
  match_id uuid references matches(id),
  user_id uuid references users(id),
  current_rating integer NOT NULL,
  is_winner boolean NOT NULL,
  score integer,
  PRIMARY KEY(match_id, user_id)
);

