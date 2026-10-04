CREATE TABLE users (
    id text PRIMARY KEY,
    name text NOT NULL CHECK (name <> '')
);

CREATE TABLE games (
    id text PRIMARY KEY,
    document jsonb NOT NULL CHECK (jsonb_typeof(document) = 'object' AND document->>'id' = id),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX games_started_players_idx ON games USING gin ((document->'players'))
    WHERE document->>'state' = 'started';
