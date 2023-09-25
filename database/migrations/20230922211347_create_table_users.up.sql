
CREATE TABLE IF NOT EXISTS roles (
    id serial PRIMARY KEY,
    name VARCHAR(20)
);
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(100),
    password VARCHAR(100),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    image_url VARCHAR(200),
    roles_id INT,
    FOREIGN KEY (roles_id) REFERENCES roles (id)
);
CREATE INDEX IF NOT EXISTS idx_role_id ON users (roles_id);
INSERT INTO roles (name) VALUES ('admin');

INSERT INTO roles (name) VALUES ('user');