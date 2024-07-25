-- name: ListRoles :many
SELECT * from role;

-- name: CreateUser :one
INSERT INTO "user"  
  (email, name, lastname, password_hash, role_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING id, name, email,lastname,role_id;

-- name: FindUserByEmail :one
SELECT  * from "user"
  where email = $1; 
