CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,                          -- Auto-incrementing ID
    first_name VARCHAR(255) NOT NULL,              -- First name
    last_name VARCHAR(255) NOT NULL,               -- Last name
    sex VARCHAR(6) CHECK (sex IN ('male','female')) NOT NULL,  -- Gender
    email VARCHAR(255) UNIQUE NOT NULL,            -- Email
    dob DATE NOT NULL,                             -- Date of Birth
    contact VARCHAR(20) NOT NULL,                  -- Contact number
    password VARCHAR(255) NOT NULL,                -- Hashed password
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),   -- Record creation time
    role VARCHAR(10) CHECK (role IN ('admin','customer','agent')) NOT NULL DEFAULT 'customer'
);
