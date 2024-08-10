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

-- name: ListUserGroups :many
SELECT * From "group" as g
  JOIN group_user as gu on g.id = gu.group_id
  WHERE gu.user_id = $1;

-- name: CreateGroup :one
INSERT INTO "group"
  (name) 
  VALUES ($1)
  RETURNING id, name;

-- name: AddUserToGroup :exec
INSERT INTO group_user
  (group_id, user_id, is_admin)
  VALUES($1, $2, $3);
