package model

import (
	"database/sql"
	"time"

	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
)

func toInt64Ptr(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

func toStringPtr(n sql.NullString) *string {
	if n.Valid {
		return &n.String
	}
	return nil
}

func toTimePtr(n sql.NullTime) *string {
	if n.Valid {
		s := n.Time.Format(time.RFC3339)
		return &s
	}
	return nil
}

func toBoolPtr(n sql.NullBool) *bool {
	if n.Valid {
		return &n.Bool
	}
	return nil
}

func AlbumFromDB(a db.Album) Album {
	return Album{
		ID:              a.ID,
		DeezerID:        toStringPtr(a.DeezerID),
		Title:           a.Title,
		NormalizedTitle: a.NormalizedTitle,
		ArtistID:        a.ArtistID,
		ReleaseDate:     toStringPtr(a.ReleaseDate),
		AlbumArtPath:    toStringPtr(a.AlbumArtPath),
		Genre:           toStringPtr(a.Genre),
		TotalTracks:     toInt64Ptr(a.TotalTracks),
		CreatedAt:       a.CreatedAt.Format(time.RFC3339),
	}
}

func ArtistFromDB(a db.Artist) Artist {
	return Artist{
		ID:             a.ID,
		DeezerID:       toStringPtr(a.DeezerID),
		Name:           a.Name,
		NormalizedName: a.NormalizedName,
		CreatedAt:      a.CreatedAt.Format(time.RFC3339),
	}
}

func DownloadHistoryFromDB(d db.DownloadHistory) DownloadHistory {
	return DownloadHistory{
		ID:           d.ID,
		UserID:       toInt64Ptr(d.UserID),
		TrackID:      toInt64Ptr(d.TrackID),
		Quality:      toInt64Ptr(d.Quality),
		Status:       toStringPtr(d.Status),
		Service:      toStringPtr(d.Service),
		StartedAt:    d.StartedAt.Format(time.RFC3339),
		CompletedAt:  toTimePtr(d.CompletedAt),
		ErrorMessage: toStringPtr(d.ErrorMessage),
	}
}

func TrackFromDB(t db.Track) Track {
	return Track{
		ID:              t.ID,
		Title:           t.Title,
		NormalizedTitle: t.NormalizedTitle,
		ArtistID:        toInt64Ptr(t.ArtistID),
		AlbumID:         toInt64Ptr(t.AlbumID),
		Duration:        toInt64Ptr(t.Duration),
		TrackNumber:     toInt64Ptr(t.TrackNumber),
		DiscNumber:      toInt64Ptr(t.DiscNumber),
		SampleRate:      toInt64Ptr(t.SampleRate),
		Bitrate:         toInt64Ptr(t.Bitrate),
		Channels:        toInt64Ptr(t.Channels),
		FilePath:        t.FilePath,
		FileSize:        toInt64Ptr(t.FileSize),
		ISRC:            toStringPtr(t.Isrc),
		CreatedAt:       t.CreatedAt.Format(time.RFC3339),
	}
}

func UserFromDB(u db.User) User {
	return User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     toStringPtr(u.Email),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		LastLogin: toTimePtr(u.LastLogin),
		IsActive:  toBoolPtr(u.IsActive),
	}
}

func UserTrackFromDB(ut db.UserTrack) UserTrack {
	return UserTrack{
		UserID:       toInt64Ptr(ut.UserID),
		TrackID:      toInt64Ptr(ut.TrackID),
		SymlinkPath:  ut.SymlinkPath,
		DownloadDate: ut.DownloadDate.Format(time.RFC3339),
	}
}

func MapToTrackPreviews(results []StreamripSearchResult) []TrackPreview {
	previews := make([]TrackPreview, 0, len(results))
	for _, r := range results {
		preview := TrackPreview{
			Title:    r.Data.Title,
			Artist:   r.Data.Performer.Name,
			Duration: r.Data.Duration,
			Image:    r.Data.Album.Image.Small,
			TrackID:  r.ID,
			Source:   r.Source,
			ISRC:     r.Data.ISRC,
		}
		previews = append(previews, preview)
	}
	return previews
}

func LimitResults[T any](items []T, max int) []T {
	if len(items) > max {
		return items[:max]
	}
	return items
}
