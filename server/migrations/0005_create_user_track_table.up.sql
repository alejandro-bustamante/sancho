-- Tabla para manejar la relación entre usuarios y tracks (descargas y symlinks)
CREATE TABLE user_track (
    user_id INTEGER,
    track_id INTEGER,
    symlink_path TEXT NOT NULL, -- ruta del symlink en la carpeta del usuario
    download_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (track_id) REFERENCES track(id),
    PRIMARY KEY (user_id, track_id)
);