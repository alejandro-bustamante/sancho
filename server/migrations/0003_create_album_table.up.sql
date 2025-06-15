-- Tabla de álbum
CREATE TABLE album (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    normalized_title TEXT NOT NULL,
    artist_id INTEGER,
    release_date TEXT,
    album_art_path TEXT,
    genre TEXT,
    year INTEGER,
    total_tracks INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES artist(id)
);