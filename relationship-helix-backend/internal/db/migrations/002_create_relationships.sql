-- Crearea tabelei pentru relații
CREATE TABLE IF NOT EXISTS relationships (
    id SERIAL PRIMARY KEY,
    user1_id INTEGER NOT NULL REFERENCES users(id),
    user2_id INTEGER NOT NULL REFERENCES users(id),
    user1_name VARCHAR(100) NOT NULL,
    user2_name VARCHAR(100) NOT NULL,
    start_date TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- O relație unică între doi utilizatori
    UNIQUE(user1_id, user2_id)
);

-- Indecși pentru performanță
CREATE INDEX idx_relationships_user1_id ON relationships(user1_id);
CREATE INDEX idx_relationships_user2_id ON relationships(user2_id);

-- Crearea tabelei pentru coduri de invitație
CREATE TABLE IF NOT EXISTS invite_codes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    code VARCHAR(20) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indecși pentru performanță
CREATE INDEX idx_invite_codes_user_id ON invite_codes(user_id);
CREATE INDEX idx_invite_codes_code ON invite_codes(code);