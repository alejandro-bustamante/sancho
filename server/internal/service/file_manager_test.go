package service_test

import (
	"context"
	_ "os"
	_ "path/filepath"
	"testing"

	"github.com/alejandro-bustamante/sancho/server/internal/config"
	model "github.com/alejandro-bustamante/sancho/server/internal/model"
	"github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/alejandro-bustamante/sancho/server/internal/testutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

type FileManagerSuite struct {
	suite.Suite
	db  repository.Database
	fs  afero.Fs
	fm  *service.FileManager
	ctx context.Context
}

func (s *FileManagerSuite) SetupTest() {
	// DB in memory
	s.db = testutils.SetupTestDB(s.T())

	// MemMapFs mocks a complete disk in RAM
	s.fs = afero.NewMemMapFs()

	s.fm = service.NewFileManager(s.db, s.fs)
	s.ctx = context.Background()

	// 4. Configurar paths globales para el test (mockear config si es necesario)
	// Nota: Como config.SanchoPath es global, asegúrate de que sea una ruta simple
	// o configúralo en el init() del test si es variable.
	// Asumiremos "/sancho" para este entorno virtual.
	config.SanchoPath = "/sancho"
}

func TestFileManagerSuite(t *testing.T) {
	suite.Run(t, new(FileManagerSuite))
}

func (s *FileManagerSuite) TestMoveTrackToLibrary() {
	// ARRANGE
	// 1. Crear datos en DB
	artist := CreateTestArtist(s.T(), s.db, "Queen")
	album := CreateTestAlbum(s.T(), s.db, "A Night at the Opera", artist.ID)
	track := CreateTestTrack(s.T(), s.db, "Bohemian Rhapsody", artist.ID, album.ID)

	// 2. Crear el archivo "físico" en el FS virtual
	// El path original en DB es "/music/Bohemian Rhapsody.flac" (del helper)
	err := s.fs.MkdirAll("/music", 0755)
	s.Require().NoError(err)
	err = afero.WriteFile(s.fs, track.FilePath, []byte("fake audio content"), 0644)
	s.Require().NoError(err)

	// ACT
	// Convertir DB Track a Model Track (necesitas un mapper o hacerlo manual)
	modelTrack := model.TrackFromDB(track) // Asumiendo que tienes este mapper disponible
	newPath, err := s.fm.MoveTrackToLibrary(s.ctx, modelTrack)

	// ASSERT
	s.Require().NoError(err)

	// 1. Verificar ruta retornada
	expectedPath := "/sancho/library/Queen/A Night at the Opera/Bohemian Rhapsody.flac"
	s.Assert().Equal(expectedPath, newPath)

	// 2. Verificar que el archivo existe en la nueva ruta
	exists, _ := afero.Exists(s.fs, newPath)
	s.Assert().True(exists, "File should exist in library")

	// 3. Verificar que ya NO existe en la vieja ruta
	existsOld, _ := afero.Exists(s.fs, track.FilePath)
	s.Assert().False(existsOld, "File should not exist in old location")

	// 4. Verificar que la DB se actualizó
	updatedTrack, _ := s.db.GetTrackByID(s.ctx, track.ID)
	s.Assert().Equal(expectedPath, updatedTrack.FilePath)
}
