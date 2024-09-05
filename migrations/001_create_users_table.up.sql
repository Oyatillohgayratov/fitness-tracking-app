CREATE TABLE IF NOT EXISTS "users" (
  "id" serial PRIMARY KEY,
  "username" varchar,
  "email" varchar,
  "password_hash" varchar,
  "profile" jsonb
);