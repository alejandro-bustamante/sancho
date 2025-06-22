-- Tabla para el historial de descargas
CREATE TABLE download_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    track_id INTEGER,
    quality INTEGER CHECK(quality IN (0, 1, 2, 3)), -- calidad de la descarga (0-3)
    status TEXT CHECK(status IN ('success', 'failed', 'pending')),
    service TEXT, -- qobuz, tidal, etc.
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    error_message TEXT,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (track_id) REFERENCES track(id)
);