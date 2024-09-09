CREATE TABLE IF NOT EXISTS workouts (
    id serial primary key,
    user_id integer not null,
    name text not null,
    description text,
    date DATE not null default current_date,
    create_at timestamptz not null default now(),
    update_at timestamptz not null default now(),
    FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
);
