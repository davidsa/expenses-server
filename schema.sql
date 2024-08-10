CREATE TABLE IF NOT EXISTS role (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS "user" (
  id SERIAL PRIMARY KEY,
  email VARCHAR(100) UNIQUE NOT NULL,
  name VARCHAR(100) NOT NULL,
  lastname VARCHAR(100) NOT NULL,
  password_hash BYTEA NOT NULL,
  role_id INT,
  CONSTRAINT fk_role
    FOREIGN KEY(role_id) 
      REFERENCES role(id)
);

CREATE TABLE IF NOT EXISTS "group" (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS group_user (
  group_id INT,
  user_id INT,
  is_admin BOOLEAN,
  PRIMARY KEY (group_id, user_id),
  CONSTRAINT fk_group FOREIGN key(group_id) REFERENCES "group"(id),
  CONSTRAINT fk_user FOREIGN key(user_id) REFERENCES "user"(id)
);

