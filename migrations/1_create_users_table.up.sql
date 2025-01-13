CREATE TABLE IF NOT EXISTS users
(
  id BIGINT GENERATED ALWAYS AS IDENTITY,
  state TEXT NOT NULL,
  email VARCHAR(255) NOT NULL,
  password TEXT NOT NULL,
  name TEXT NOT NULL,
  lastname TEXT NOT NULL,
  avatar_url TEXT DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT users_pkey PRIMARY KEY (id),
  CONSTRAINT users_state_check CHECK (state IN ('registered', 'verified', 'banned', 'deleted')),
  CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE OR REPLACE RULE delete_users_rule AS
    ON DELETE TO users
    DO INSTEAD (
        UPDATE users
        SET deleted_at = now()
        WHERE id = OLD.id
    );
