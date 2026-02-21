package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/alejandro-bustamante/sancho/server/internal/testutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	tag "go.senan.xyz/taglib"
)

// MockMetadataExtractor implementa la interfaz MetadataExtractor
type MockMetadataExtractor struct {
	mock.Mock
}

func (m *MockMetadataExtractor) ReadTags(path string) (map[string][]string, error) {
	args := m.Called(path)
	return args.Get(0).(map[string][]string), args.Error(1)
}

func (m *MockMetadataExtractor) ReadProperties(path string) (tag.Properties, error) {
	args := m.Called(path)
	return args.Get(0).(tag.Properties), args.Error(1)
}

type IndexerSuite struct {
	suite.Suite
	db       repository.Database
	fs       afero.Fs
	indexer  *service.Indexer
	mockMeta *MockMetadataExtractor
}

func (s *IndexerSuite) SetupTest() {
	s.db = testutils.SetupTestDB(s.T())

	// We need to work with a real dir for the symlinks
	// The simpler afero MemMapFs does not implement symlinks
	baseDir := s.T().TempDir()
	s.fs = afero.NewBasePathFs(afero.NewOsFs(), baseDir)

	s.mockMeta = new(MockMetadataExtractor)

	// FileManager también necesita el FS
	fm := service.NewFileManager(s.db, s.fs)

	// Inyectamos todo en el indexer
	s.indexer = service.NewIndexer(s.db, fm, s.fs, s.mockMeta)
}

func TestIndexerSuite(t *testing.T) {
	suite.Run(t, new(IndexerSuite))
}

func (s *IndexerSuite) TestIndexFolder() {
	// 1. PREPARAR: Crear el usuario en la DB (necesario para saveTransferHistory)
	_, err := s.db.InsertUser(context.Background(), repository.InsertUserParams{
		Username:     "admin",
		Email:        sql.NullString{String: "admin@example.com", Valid: true},
		PasswordHash: "hash",
	})
	s.Require().NoError(err)

	// 2. Estructura de archivos
	root := "/sancho/admin_library/Test Artist/Test Album"
	s.fs.MkdirAll(root, 0755)
	songPath := root + "/song.mp3"
	afero.WriteFile(s.fs, songPath, []byte("mp3 data"), 0644)

	// 3. Configurar Mock de Metadatos
	tags := map[string][]string{
		"TITLE":  {"Test Song"},
		"ARTIST": {"Test Artist"},
		"ALBUM":  {"Test Album"},
		"ISRC":   {"US1234567890"},
	}
	props := tag.Properties{Length: time.Second * 180, Bitrate: 320}

	s.mockMeta.On("ReadTags", songPath).Return(tags, nil)
	s.mockMeta.On("ReadProperties", songPath).Return(props, nil)

	// ACT
	err = s.indexer.IndexFolder(context.Background(), root, "admin", "local", 320)

	// ASSERT
	s.Require().NoError(err, "El indexer falló. Revisa los logs para ver errores de Deezer o DB")

	// Verificar Artista
	artist, err := s.db.GetArtistByNormalizedName(context.Background(), "test artist")
	s.Require().NoError(err)
	s.Assert().Equal("Test Artist", artist.Name)

	// Verificar Track
	results, err := s.db.SearchTracksByTitle(context.Background(), sql.NullString{String: "Test Song", Valid: true})
	s.Require().NoError(err)
	s.Require().NotEmpty(results)
	s.Assert().Equal("US1234567890", results[0].Isrc.String)
}
