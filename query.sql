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
