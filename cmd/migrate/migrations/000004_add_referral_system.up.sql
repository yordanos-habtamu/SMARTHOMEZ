-- Create referral_codes table
CREATE TABLE IF NOT EXISTS referral_codes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    qr_code_url VARCHAR(500),
    visit_count INT NOT NULL DEFAULT 0,
    signup_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_referral_codes_user_id ON referral_codes(user_id);
CREATE INDEX idx_referral_codes_short_code ON referral_codes(short_code);

-- Create referral_analytics table for tracking
CREATE TABLE IF NOT EXISTS referral_analytics (
    id SERIAL PRIMARY KEY,
    referral_code_id INT NOT NULL,
    event_type VARCHAR(20) CHECK(event_type IN ('visit', 'signup')) NOT NULL,
    ip_address VARCHAR(50),
    user_agent TEXT,
    converted_user_id INT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (referral_code_id) REFERENCES referral_codes(id) ON DELETE CASCADE,
    FOREIGN KEY (converted_user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes for analytics queries
CREATE INDEX idx_referral_analytics_code_id ON referral_analytics(referral_code_id);
CREATE INDEX idx_referral_analytics_event_type ON referral_analytics(event_type);
CREATE INDEX idx_referral_analytics_created_at ON referral_analytics(created_at);
