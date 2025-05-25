CREATE TABLE IF NOT EXISTS block_areas (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    radius FLOAT NOT NULL,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    altitude FLOAT NOT NULL,
    state VARCHAR(20) NOT NULL DEFAULT 'active',
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_block_areas_user_id ON block_areas(user_id);
CREATE INDEX idx_block_areas_state ON block_areas(state);
CREATE INDEX idx_block_areas_expires_at ON block_areas(expires_at); 