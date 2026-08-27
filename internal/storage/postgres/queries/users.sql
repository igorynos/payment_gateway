-- name: CreateUser :one
  INSERT INTO users (
      username,
      password_hash,
      email,
      role
  )
  VALUES (
      sqlc.arg(username),
      sqlc.arg(password_hash),
      sqlc.arg(email),
      sqlc.arg(role)
  )
  RETURNING *;

-- name: GetUserByID :one
  SELECT *
  FROM users
  WHERE id = $1;

-- name: ListUsers :one
  SELECT id
  FROM users;

-- name: DeleteUser :exec
  DELETE FROM users
  WHERE id = $1;
