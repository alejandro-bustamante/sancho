package service_test

import (
	"context"
	"testing"

	"github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/alejandro-bustamante/sancho/server/internal/testutils"
	"github.com/stretchr/testify/suite"
)

type LibraryServiceSuite struct {
	suite.Suite
	db      repository.Database
	service *service.LibraryService
	ctx     context.Context
}

func (s *LibraryServiceSuite) SetupTest() {
	s.db = testutils.SetupTestDB(s.T())
	s.service = service.NewLibraryService(s.db)
	s.ctx = context.Background()
}

func TestLibraryServiceSuite(t *testing.T) {
	suite.Run(t, new(LibraryServiceSuite))
}

func (s *LibraryServiceSuite) TestGetUserTracks() {
	// Arrange: Prepare complex data hierarchy
	userService := service.NewUserService(s.db)
	user, err := userService.RegisterUser(s.ctx, "music_lover", "pass", "test@test.com")
	s.Require().NoError(err)

	artist := CreateTestArtist(s.T(), s.db, "Daft Punk")
	album := CreateTestAlbum(s.T(), s.db, "Discovery", artist.ID)

	track1 := CreateTestTrack(s.T(), s.db, "One More Time", artist.ID, album.ID)
	_ = CreateTestTrack(s.T(), s.db, "Aerodynamic", artist.ID, album.ID)
	track3 := CreateTestTrack(s.T(), s.db, "Harder, Better", artist.ID, album.ID)

	// Link only 2 tracks to the user
	LinkTrackToUser(s.T(), s.db, user.ID, track1.ID)
	LinkTrackToUser(s.T(), s.db, user.ID, track3.ID)

	result, err := s.service.GetUserTracks(s.ctx, "music_lover")

	s.Require().NoError(err)
	s.Assert().Len(result, 2, "User should have exactly 2 tracks in library")

	// Verify we got the correct tracks
	titles := []string{result[0].Title, result[1].Title}
	s.Assert().Contains(titles, "One More Time")
	s.Assert().Contains(titles, "Harder, Better")
	s.Assert().NotContains(titles, "Aerodynamic")
}

func (s *LibraryServiceSuite) TestSearchTracks() {
	// Arrange: Seed DB with searchable content
	artist := CreateTestArtist(s.T(), s.db, "The Beatles")
	album := CreateTestAlbum(s.T(), s.db, "Abbey Road", artist.ID)

	CreateTestTrack(s.T(), s.db, "Come Together", artist.ID, album.ID)
	CreateTestTrack(s.T(), s.db, "Something", artist.ID, album.ID)
	CreateTestTrack(s.T(), s.db, "Here Comes The Sun", artist.ID, album.ID)

	// Exact match (partial)
	results, err := s.service.SearchTracks(s.ctx, "Come")
	s.Require().NoError(err)
	s.Assert().Len(results, 2, "Should find 'Come Together' and 'Here Comes The Sun'")

	// Case insensitive
	resultsLower, err := s.service.SearchTracks(s.ctx, "something")
	s.Require().NoError(err)
	s.Assert().Len(resultsLower, 1)
	s.Assert().Equal("Something", resultsLower[0].Title)

	// No results
	resultsEmpty, err := s.service.SearchTracks(s.ctx, "Nirvana")
	s.Require().NoError(err)
	s.Assert().Empty(resultsEmpty)
}

func (s *LibraryServiceSuite) TestGetTrackByID() {
	artist := CreateTestArtist(s.T(), s.db, "Solo Artist")
	album := CreateTestAlbum(s.T(), s.db, "Solo Album", artist.ID)
	track := CreateTestTrack(s.T(), s.db, "Lonely Track", artist.ID, album.ID)

	foundTrack, err := s.service.GetTrackByID(s.ctx, track.ID)

	s.Require().NoError(err)
	s.Assert().Equal(track.Title, foundTrack.Title)
	s.Assert().Equal(track.FilePath, foundTrack.FilePath)
}
