CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid NOT NULL,
  "username" text NOT NULL,
  "fullname" text NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
    PRIMARY KEY("id")
);
