-- Tabla de track (canción)
CREATE TABLE track (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    normalized_title TEXT NOT NULL,
    artist_id INTEGER,
    album_id INTEGER,
    duration INTEGER, -- en milisegundos
    track_number INTEGER,
    disc_number INTEGER DEFAULT 1,
    sample_rate INTEGER,
    -- bit_depth INTEGER,
    bitrate INTEGER,
    channels INTEGER,
    -- codec TEXT,
    file_path TEXT NOT NULL, -- ruta en la carpeta principal
    file_size INTEGER,
    isrc TEXT,
    composer TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES artist(id),
    FOREIGN KEY (album_id) REFERENCES album(id)
);