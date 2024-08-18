CREATE TABLE sessions (
      guid TEXT PRIMARY KEY,
      refresh_token BYTEA,
      live_time TIMESTAMP,
      client_ip TEXT
);