package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Estructuras para el JSON
type TestData struct {
	Users   []UserData   `json:"users"`
	Artists []ArtistData `json:"artists"`
}

type UserData struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Tracks   []string `json:"tracks"` // ISRCs de las canciones que tiene
}

type ArtistData struct {
	Name     string      `json:"name"`
	Albums   []AlbumData `json:"albums"`
	DeezerID string      `json:"deezer_id"`
}

type AlbumData struct {
	Title       string      `json:"title"`
	ReleaseDate string      `json:"release_date"`
	Genre       string      `json:"genre"`
	DeezerID    string      `json:"deezer_id"`
	Tracks      []TrackData `json:"tracks"`
}

type TrackData struct {
	Title       string `json:"title"`
	Duration    int    `json:"duration"` // milisegundos
	TrackNumber int    `json:"track_number"`
	ISRC        string `json:"isrc"`
	Composer    string `json:"composer"`
}

const (
	testEnvDir   = "./test_env"
	dbPath       = "./test_env/sancho_test.db"
	libraryPath  = "./test_env/library"
	testDataJSON = "./test_data.json"
)

func main() {
	log.Println("Iniciando creación del entorno de pruebas...")

	// Limpiar entorno previo
	cleanupEnvironment()

	// Leer datos de prueba
	testData, err := loadTestData()
	if err != nil {
		log.Fatalf("Error cargando datos de prueba: %v", err)
	}

	// Crear base de datos
	db, err := createDatabase()
	if err != nil {
		log.Fatalf("Error creando base de datos: %v", err)
	}
	defer db.Close()

	// Poblar base de datos
	if err := populateDatabase(db, testData); err != nil {
		log.Fatalf("Error poblando base de datos: %v", err)
	}

	// Crear estructura de archivos
	if err := createFileStructure(db, testData); err != nil {
		log.Fatalf("Error creando estructura de archivos: %v", err)
	}

	log.Println("✓ Entorno de pruebas creado exitosamente")
}

func cleanupEnvironment() {
	os.RemoveAll(testEnvDir)
	os.MkdirAll(testEnvDir, 0755)
}

func loadTestData() (*TestData, error) {
	// Si no existe el archivo JSON, crear uno por defecto
	if _, err := os.Stat(testDataJSON); os.IsNotExist(err) {
		if err := createDefaultTestData(); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(testDataJSON)
	if err != nil {
		return nil, err
	}

	var testData TestData
	if err := json.Unmarshal(data, &testData); err != nil {
		return nil, err
	}

	return &testData, nil
}

func createDefaultTestData() error {
	defaultData := TestData{
		Users: []UserData{
			{
				Username: "alejandro",
				Email:    "alejandro@test.com",
				Password: "password123",
				Tracks:   []string{"USUM71703692", "DEA621700390", "BRA340800123"},
			},
			{
				Username: "maria",
				Email:    "maria@test.com",
				Password: "password123",
				Tracks:   []string{"USUM71703692", "GBAYE0601477", "BRA340800123"},
			},
		},
		Artists: []ArtistData{
			{
				Name:     "Daft Punk",
				DeezerID: "27",
				Albums: []AlbumData{
					{
						Title:       "Random Access Memories",
						ReleaseDate: "2013-05-17",
						Genre:       "Electronic",
						DeezerID:    "6575789",
						Tracks: []TrackData{
							{
								Title:       "Get Lucky (feat. Pharrell Williams)",
								Duration:    369000,
								TrackNumber: 8,
								ISRC:        "USUM71703692",
								Composer:    "Thomas Bangalter, Guy-Manuel de Homem-Christo",
							},
						},
					},
				},
			},
			{
				Name:     "Kraftwerk",
				DeezerID: "524",
				Albums: []AlbumData{
					{
						Title:       "Die Mensch-Maschine",
						ReleaseDate: "1978-05-01",
						Genre:       "Electronic",
						DeezerID:    "302127",
						Tracks: []TrackData{
							{
								Title:       "Das Model",
								Duration:    218000,
								TrackNumber: 5,
								ISRC:        "DEA621700390",
								Composer:    "Ralf Hütter, Karl Bartos",
							},
						},
					},
				},
			},
			{
				Name:     "Caetano Veloso",
				DeezerID: "5290",
				Albums: []AlbumData{
					{
						Title:       "Transa",
						ReleaseDate: "1972-01-01",
						Genre:       "MPB",
						DeezerID:    "12345678",
						Tracks: []TrackData{
							{
								Title:       "Você Não Entende Nada",
								Duration:    245000,
								TrackNumber: 3,
								ISRC:        "BRA340800123",
								Composer:    "Caetano Veloso",
							},
						},
					},
				},
			},
			{
				Name:     "Adele",
				DeezerID: "75491",
				Albums: []AlbumData{
					{
						Title:       "21",
						ReleaseDate: "2011-01-24",
						Genre:       "Pop, Soul",
						DeezerID:    "1223043",
						Tracks: []TrackData{
							{
								Title:       "Rolling in the Deep",
								Duration:    228000,
								TrackNumber: 1,
								ISRC:        "GBAYE0601477",
								Composer:    "Adele Adkins, Paul Epworth",
							},
						},
					},
				},
			},
		},
	}

	data, err := json.MarshalIndent(defaultData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(testDataJSON, data, 0644)
}

func createDatabase() (*sql.DB, error) {
	// Asegurar que el directorio existe
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE schema_migrations (version INTEGER, dirty INTEGER);
	CREATE UNIQUE INDEX version_unique ON schema_migrations (version);
	
	CREATE TABLE user (
		id INTEGER PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		email TEXT UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		last_login TIMESTAMP,
		is_active INTEGER DEFAULT 1
	);
	
	CREATE TABLE artist (
		id INTEGER PRIMARY KEY,
		deezer_id TEXT UNIQUE,
		name TEXT NOT NULL,
		normalized_name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	);
	
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
	
	CREATE TABLE track (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		normalized_title TEXT NOT NULL,
		artist_id INTEGER,
		album_id INTEGER,
		duration INTEGER,
		track_number INTEGER,
		disc_number INTEGER DEFAULT 1,
		sample_rate INTEGER,
		bitrate INTEGER,
		channels INTEGER,
		file_path TEXT NOT NULL,
		file_size INTEGER,
		isrc TEXT,
		composer TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (artist_id) REFERENCES artist(id),
		FOREIGN KEY (album_id) REFERENCES album(id)
	);
	
	CREATE TABLE user_track (
		user_id INTEGER,
		track_id INTEGER,
		symlink_path TEXT NOT NULL,
		linked_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES user(id),
		FOREIGN KEY (track_id) REFERENCES track(id),
		PRIMARY KEY (user_id, track_id)
	);
	
	CREATE TABLE download_history (
		id TEXT PRIMARY KEY,
		user_id INTEGER,
		track_id INTEGER,
		quality INTEGER CHECK(quality IN (0, 1, 2, 3)),
		status TEXT CHECK(status IN ('success', 'downloading', 'indexing', 'failed', 'canceled', 'transfered')),
		service TEXT,
		started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		completed_at TIMESTAMP,
		error_message TEXT,
		FOREIGN KEY (user_id) REFERENCES user(id),
		FOREIGN KEY (track_id) REFERENCES track(id)
	);
	
	CREATE INDEX idx_track_title ON track(normalized_title);
	CREATE INDEX idx_album_title ON album(normalized_title);
	CREATE UNIQUE INDEX idx_album_deezer_id ON album(deezer_id);
	CREATE INDEX idx_artist_name ON artist(normalized_name);
	CREATE UNIQUE INDEX idx_artist_deezer_id ON artist(deezer_id);
	CREATE INDEX idx_user_track_user_id ON user_track(user_id);
	CREATE INDEX idx_track_isrc ON track(isrc);
	
	INSERT INTO schema_migrations (version, dirty) VALUES (7, 0);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func populateDatabase(db *sql.DB, testData *TestData) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insertar usuarios
	userIDs := make(map[string]int64)
	for _, user := range testData.Users {
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		result, err := tx.Exec(`
			INSERT INTO user (username, password_hash, email, is_active)
			VALUES (?, ?, ?, 1)
		`, user.Username, string(hash), user.Email)
		if err != nil {
			return err
		}
		userIDs[user.Username], _ = result.LastInsertId()
	}

	// Insertar artistas, álbumes y canciones
	trackISRCtoID := make(map[string]int64)

	for _, artist := range testData.Artists {
		// Insertar artista
		artistResult, err := tx.Exec(`
			INSERT INTO artist (deezer_id, name, normalized_name)
			VALUES (?, ?, ?)
		`, artist.DeezerID, artist.Name, normalize(artist.Name))
		if err != nil {
			return err
		}
		artistID, _ := artistResult.LastInsertId()

		for _, album := range artist.Albums {
			// Insertar álbum
			albumResult, err := tx.Exec(`
				INSERT INTO album (deezer_id, title, normalized_title, artist_id, release_date, genre, total_tracks)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, album.DeezerID, album.Title, normalize(album.Title), artistID, album.ReleaseDate, album.Genre, len(album.Tracks))
			if err != nil {
				return err
			}
			albumID, _ := albumResult.LastInsertId()

			for _, track := range album.Tracks {
				// Construir ruta del archivo
				fileName := fmt.Sprintf("%02d. %s - %s.flac", track.TrackNumber, track.Title, artist.Name)

				// Find the executable path to put the full path in the test database
				executable, _ := os.Executable()
				testDir := filepath.Dir(executable)
				filePath := filepath.Join(testDir, "/test_env/library", artist.Name, album.Title, fileName)

				// Insertar canción
				trackResult, err := tx.Exec(`
					INSERT INTO track (title, normalized_title, artist_id, album_id, duration, track_number, 
						disc_number, sample_rate, bitrate, channels, file_path, file_size, isrc, composer)
					VALUES (?, ?, ?, ?, ?, ?, 1, 44100, 988, 2, ?, 37955782, ?, ?)
				`, track.Title, normalize(track.Title), artistID, albumID, track.Duration,
					track.TrackNumber, filePath, track.ISRC, track.Composer)
				if err != nil {
					return err
				}
				trackID, _ := trackResult.LastInsertId()
				trackISRCtoID[track.ISRC] = trackID
			}
		}
	}

	// Insertar user_track y download_history
	for _, user := range testData.Users {
		userID := userIDs[user.Username]
		for _, isrc := range user.Tracks {
			trackID, exists := trackISRCtoID[isrc]
			if !exists {
				continue
			}

			// Obtener información de la canción para construir el symlink
			var artistName, albumTitle, trackTitle string
			var trackNumber int
			err := tx.QueryRow(`
				SELECT a.name, al.title, t.title, t.track_number
				FROM track t
				JOIN artist a ON t.artist_id = a.id
				JOIN album al ON t.album_id = al.id
				WHERE t.id = ?
			`, trackID).Scan(&artistName, &albumTitle, &trackTitle, &trackNumber)
			if err != nil {
				continue
			}

			fileName := fmt.Sprintf("%02d. %s - %s.flac", trackNumber, trackTitle, artistName)
			symlinkPath := filepath.Join("../../../library", artistName, albumTitle, fileName)

			// Insertar user_track
			_, err = tx.Exec(`
				INSERT INTO user_track (user_id, track_id, symlink_path)
				VALUES (?, ?, ?)
			`, userID, trackID, symlinkPath)
			if err != nil {
				return err
			}

			// Insertar download_history
			downloadID := generateUUID()
			_, err = tx.Exec(`
				INSERT INTO download_history (id, user_id, track_id, quality, status, service, completed_at)
				VALUES (?, ?, ?, 2, 'transfered', 'qobuz', CURRENT_TIMESTAMP)
			`, downloadID, userID, trackID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func createFileStructure(db *sql.DB, testData *TestData) error {
	// Crear librería principal
	if err := os.MkdirAll(libraryPath, 0755); err != nil {
		return err
	}

	// Crear archivos de canciones
	for _, artist := range testData.Artists {
		for _, album := range artist.Albums {
			albumPath := filepath.Join(libraryPath, artist.Name, album.Title)
			if err := os.MkdirAll(albumPath, 0755); err != nil {
				return err
			}

			for _, track := range album.Tracks {
				fileName := fmt.Sprintf("%02d. %s - %s.flac", track.TrackNumber, track.Title, artist.Name)
				filePath := filepath.Join(albumPath, fileName)
				// Crear archivo dummy
				if err := os.WriteFile(filePath, []byte("dummy audio data"), 0644); err != nil {
					return err
				}
			}
		}
	}

	// Crear librerías de usuarios con symlinks
	for _, user := range testData.Users {
		userLibPath := filepath.Join(testEnvDir, fmt.Sprintf("%s_library", user.Username))
		if err := os.MkdirAll(userLibPath, 0755); err != nil {
			return err
		}

		for _, isrc := range user.Tracks {
			var artistName, albumTitle, trackTitle string
			var trackNumber int
			err := db.QueryRow(`
				SELECT a.name, al.title, t.title, t.track_number
				FROM track t
				JOIN artist a ON t.artist_id = a.id
				JOIN album al ON t.album_id = al.id
				WHERE t.isrc = ?
			`, isrc).Scan(&artistName, &albumTitle, &trackTitle, &trackNumber)
			if err != nil {
				continue
			}

			// Crear estructura de directorios
			userAlbumPath := filepath.Join(userLibPath, artistName, albumTitle)
			if err := os.MkdirAll(userAlbumPath, 0755); err != nil {
				return err
			}

			// Crear symlink
			fileName := fmt.Sprintf("%02d. %s - %s.flac", trackNumber, trackTitle, artistName)
			symlinkPath := filepath.Join(userAlbumPath, fileName)
			targetPath := filepath.Join("../../../library", artistName, albumTitle, fileName)

			if err := os.Symlink(targetPath, symlinkPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
