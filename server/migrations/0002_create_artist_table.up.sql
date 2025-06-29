-- Tabla de artista
CREATE TABLE artist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    deezer_id TEXT UNIQUE,
    name TEXT NOT NULL,
    normalized_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);