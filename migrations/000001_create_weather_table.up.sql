BEGIN;

CREATE TABLE IF NOT EXISTS weather (
  id serial PRIMARY KEY,
  city text NOT NULL,
  dt timestamp NOT NULL,
  temperature numeric NOT NULL
);

COMMIT;