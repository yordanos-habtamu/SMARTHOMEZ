CREATE TABLE IF NOT EXISTS House (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY, -- Auto-incrementing ID
    `type` ENUM('apartment', 'villa', 'house','luxury','condo') NOT NULL, -- Type of house
    `address` VARCHAR(255) NOT NULL, -- Address of the house
    `price` DECIMAL(10, 2) NOT NULL, -- Price of the house
    `num_bedrooms` INT NOT NULL, -- Number of bedrooms
    `num_bathrooms` INT NOT NULL, -- Number of bathrooms
    `area_sq_ft` DECIMAL(10,2) NOT NULL, -- Area in square feet
    `description` TEXT, -- Description of the house
    `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(), -- Record creation time
    `agent_id` INT UNSIGNED DEFAULT NULL,
    `img_url` VARCHAR(255) NOT NULL,
    `is_sold` BOOLEAN NOT NULL DEFAULT FALSE, -- Sale status
    FOREIGN KEY (agent_id) REFERENCES Users(id) ON DELETE SET NULL
);
