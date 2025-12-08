CREATE TABLE IF NOT EXISTS Users (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY, -- Auto-incrementing ID
    `first_name` VARCHAR(255) NOT NULL,           -- First name of the user
    `last_name` VARCHAR(255) NOT NULL,            -- Last name of the user
    `sex` ENUM('male', 'female') NOT NULL,        -- Gender/Sex
    `email` VARCHAR(255) UNIQUE NOT NULL,         -- Email (must be unique)
    `dob` DATE NOT NULL,                           -- Date of Birth
    `contact` VARCHAR(20) NOT NULL,          -- Contact number
    `password` VARCHAR(255) NOT NULL,             -- Hashed password
    `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),-- Record creation time
    role ENUM('admin', 'customer', 'agent') NOT NULL DEFAULT 'customer'
);
