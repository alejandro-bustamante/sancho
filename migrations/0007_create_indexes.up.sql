-- Índices para optimizar consultas comunes
CREATE INDEX idx_track_title ON track(normalized_title);
CREATE INDEX idx_album_title ON album(normalized_title);
CREATE UNIQUE INDEX idx_album_deezer_id ON album(deezer_id);
CREATE INDEX idx_artist_name ON artist(normalized_name);
CREATE UNIQUE INDEX idx_artist_deezer_id ON artist(deezer_id);
CREATE INDEX idx_user_track_user_id ON user_track(user_id);
CREATE INDEX idx_track_isrc ON track(isrc);


