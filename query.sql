-- name: GetUser :one 
select * from users
where id = $1 limit 1;

-- name: ListUser :many
select id, username, email, profile
from users
order by username;

-- name: CreateUser :one
insert into users (username, password_hash, email, profile)
values ($1, $2, $3, $4)
returning *;

-- name: UpdateUser :exec
update users
set username = $2, email = $3, profile = $4
where id = $1;

-- name: DeleteUser :exec
delete from users
where id = $1;

-- name: GetUserByEmail :one
select * from users
where email = $1 limit 1;

-- name: SavePasswordResetToken :exec
insert into password_reset_tokens (user_id, token, expiration)
values ($1, $2, $3);

-- name: GetPasswordResetToken :one
select * from password_reset_tokens
where token = $1 limit 1;

-- name: UpdatePassword :exec
update users
set password_hash = $2
where id = $1;


-- name: CreateWorkout :one
insert into workouts (user_id, name, description,date)
values ($1, $2, $3, $4)
returning id, user_id, name , description, date, create_at, update_at;

-- name: GetWorkoutsByUserID :many
select id, user_id, name, description, date, create_at, update_at
from workouts
where user_id = $1;

-- name: GetWorkoutByUserID :one
select id, user_id, name, description, date, create_at, update_at
from workouts
where id = $1 and user_id = $2;

-- name: UpdateWorkout :exec
update workouts
set name = $3, description = $4, date = $5, update_at = now()
where id = $1 and user_id = $2;

-- name: DeleteWorkout :exec
delete from workouts
where id = $1 and user_id = $2;
