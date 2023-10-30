CREATE TABLE IF NOT EXISTS posts (
    id serial PRIMARY KEY,
    name VARCHAR(40) NOT NULL,
    content TEXT,
    image_url VARCHAR(200),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);