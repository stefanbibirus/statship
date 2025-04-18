-- Crearea tabelei pentru utilizatori
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indecși pentru performanță
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

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
CREATE INDEX IF NOT EXISTS idx_relationships_user1_id ON relationships(user1_id);
CREATE INDEX IF NOT EXISTS idx_relationships_user2_id ON relationships(user2_id);

-- Crearea tabelei pentru coduri de invitație
CREATE TABLE IF NOT EXISTS invite_codes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    code VARCHAR(20) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indecși pentru performanță
CREATE INDEX IF NOT EXISTS idx_invite_codes_user_id ON invite_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_invite_codes_code ON invite_codes(code);

-- Crearea tabelei pentru pozițiile curbelor
CREATE TABLE IF NOT EXISTS curve_positions (
    id SERIAL PRIMARY KEY,
    relationship_id INTEGER NOT NULL REFERENCES relationships(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Un utilizator are o singură poziție pentru o relație
    UNIQUE(relationship_id, user_id)
);

-- Indecși pentru performanță
CREATE INDEX IF NOT EXISTS idx_curve_positions_relationship_id ON curve_positions(relationship_id);
CREATE INDEX IF NOT EXISTS idx_curve_positions_user_id ON curve_positions(user_id);