CREATE TABLE IF NOT EXISTS house (
    id SERIAL PRIMARY KEY,
    category VARCHAR(20) CHECK(category IN('apartment', 'villa', 'house','luxury','condo')) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    num_bedrooms INT NOT NULL,
    num_bathrooms INT NOT NULL,
    area_sq_ft DECIMAL(10,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    agent_id INT DEFAULT NULL,
    img_url VARCHAR(255) NOT NULL,
    is_sold BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (agent_id) REFERENCES users(id) ON DELETE SET NULL
);
