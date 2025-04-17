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
CREATE INDEX idx_curve_positions_relationship_id ON curve_positions(relationship_id);
CREATE INDEX idx_curve_positions_user_id ON curve_positions(user_id);