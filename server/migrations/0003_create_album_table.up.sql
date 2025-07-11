-- Tabla de álbum
CREATE TABLE album (
    id INTEGER PRIMARY KEY,
    deezer_id TEXT UNIQUE,
    title TEXT NOT NULL,
    normalized_title TEXT NOT NULL,
    artist_id INTEGER NOT NULL,
    release_date TEXT,
    album_art_path TEXT,
    genre TEXT,
    total_tracks INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES artist(id)
);