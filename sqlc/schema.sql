-- noinspection SqlNoDataSourceInspectionForFile

CREATE TYPE rating_system_enum AS ENUM ('glicko', 'elo');

CREATE TABLE users (
  id   uuid NOT NULL PRIMARY KEY,
  name text NOT NULL
);

CREATE TABLE games (
  id uuid NOT NULL PRIMARY KEY,
  name text NOT NULL,
  rating_system rating_system_enum NOT NULL
);

CREATE TABLE rating (
  user_id uuid references users(id),
  game_id uuid references games(id),
  rating int NOT NULL,
  PRIMARY KEY(user_id, game_id)
);

CREATE TABLE matches (
  id uuid NOT NULL PRIMARY KEY,
  game_id uuid NOT NULL references games(id),
  finished boolean NOT NULL default false,
  ratings_updated boolean NOT NULL default false,
  happened_at timestamp with time zone NOT NULL
);

CREATE TABLE match_player (
  match_id uuid references matches(id),
  user_id uuid references users(id),
  score integer,
  rating_before integer,
  rating_change integer,
  PRIMARY KEY(match_id, user_id)
);
