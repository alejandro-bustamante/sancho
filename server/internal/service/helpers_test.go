package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/stretchr/testify/require"
)

// CreateTestArtist inserts a dummy artist into the DB.
func CreateTestArtist(t *testing.T, db repository.Database, name string) repository.Artist {
	params := repository.InsertArtistParams{
		Name:           name,
		NormalizedName: name, // In real app, run through normalizer
		DeezerID:       sql.NullString{String: "123", Valid: true},
	}
	artist, err := db.InsertArtist(context.Background(), params)
	require.NoError(t, err)
	return artist
}

// CreateTestAlbum inserts a dummy album linked to an artist.
func CreateTestAlbum(t *testing.T, db repository.Database, title string, artistID int64) repository.Album {
	params := repository.InsertAlbumParams{
		Title:           title,
		NormalizedTitle: title,
		ArtistID:        artistID,
		DeezerID:        sql.NullString{String: "456", Valid: true},
		ReleaseDate:     sql.NullString{String: "2023-01-01", Valid: true},
		TotalTracks:     sql.NullInt64{Int64: 10, Valid: true},
	}
	album, err := db.InsertAlbum(context.Background(), params)
	require.NoError(t, err)
	return album
}

// CreateTestTrack inserts a dummy track linked to an album and artist.
func CreateTestTrack(t *testing.T, db repository.Database, title string, artistID, albumID int64) repository.Track {
	params := repository.InsertTrackParams{
		Title:           title,
		NormalizedTitle: title,
		ArtistID:        sql.NullInt64{Int64: artistID, Valid: true},
		AlbumID:         sql.NullInt64{Int64: albumID, Valid: true},
		FilePath:        fmt.Sprintf("/music/%s.flac", title),
		Duration:        sql.NullInt64{Int64: 300, Valid: true},
		TrackNumber:     sql.NullInt64{Int64: 1, Valid: true},
		Isrc:            sql.NullString{String: "ISRC-" + title, Valid: true},
	}
	track, err := db.InsertTrack(context.Background(), params)
	require.NoError(t, err)
	return track
}

// LinkTrackToUser connects a track to a user (simulating adding to library).
func LinkTrackToUser(t *testing.T, db repository.Database, userID, trackID int64) {
	params := repository.AddTrackToUserParams{
		UserID:      sql.NullInt64{Int64: userID, Valid: true},
		TrackID:     sql.NullInt64{Int64: trackID, Valid: true},
		SymlinkPath: "/home/user/music/link.flac",
	}
	err := db.AddTrackToUser(context.Background(), params)
	require.NoError(t, err)
}
