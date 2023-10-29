
CREATE TABLE IF NOT EXISTS cards (
    id serial PRIMARY KEY,
    name VARCHAR(40) NOT NULL,
    content TEXT,
    position INT,
    type SMALLINT,
    parent_id INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES cards(id)
);