--moncef:update_users_table
ALTER TABLE users
  RENAME COLUMN name TO first_name;

ALTER TABLE users
  ADD COLUMN last_name VARCHAR(255),
  ADD COLUMN login VARCHAR(255) UNIQUE,
  ADD COLUMN password VARCHAR(255);
