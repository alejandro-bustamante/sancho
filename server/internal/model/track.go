package model

type Track struct {
	Title string
	Duration int

}

// CREATE TABLE tracks (
//     id INTEGER PRIMARY KEY,
//     provider_id INTEGER,
//     provider_track_id TEXT, -- ID único de la canción en el proveedor
//     title TEXT NOT NULL,
//     duration INTEGER, -- en segundos
//     release_date TEXT,
//     isrc TEXT, -- código ISRC internacional
//     file_path TEXT NOT NULL, -- ruta al archivo en la carpeta general
//     file_size INTEGER,
//     sample_rate INTEGER,
//     bit_depth INTEGER,
//     channels INTEGER,
//     format TEXT, -- FLAC, MP3, etc.
//     checksum TEXT, -- hash para verificar integridad
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     FOREIGN KEY (provider_id) REFERENCES providers(id),
//     UNIQUE (provider_id, provider_track_id)
// );