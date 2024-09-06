CREATE TABLE IF NOT EXISTS "users" (
  "id" serial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password_hash" varchar NOT NULL,
  "profile" jsonb DEFAULT '{}'
);

CREATE TABLE "workouts" (
  "id" serial PRIMARY KEY,
  "user_id" int,
  "date" varchar
);

CREATE TABLE "exercises" (
  "id" serial PRIMARY KEY,
  "workout_id" int,
  "name" varchar
);

CREATE TABLE "sets" (
  "id" serial PRIMARY KEY,
  "exercise_id" int,
  "repetitions" int,
  "weight" float
);

CREATE TABLE "images" (
  "id" serial PRIMARY KEY,
  "user_id" int,
  "url" varchar
);

ALTER TABLE "workouts" 
ADD FOREIGN KEY ("user_id")
REFERENCES "users" ("id");

ALTER TABLE "exercises" 
ADD FOREIGN KEY ("workout_id") 
REFERENCES "workouts" ("id");

ALTER TABLE "sets" 
ADD FOREIGN KEY ("exercise_id") 
REFERENCES "exercises" ("id");

ALTER TABLE "images" 
ADD FOREIGN KEY ("user_id") 
REFERENCES "users" ("id");
