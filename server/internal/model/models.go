package model

type Album struct {
	ID              int64   `json:"id"`
	Title           string  `json:"title"`
	NormalizedTitle string  `json:"normalized_title"`
	ArtistID        *int64  `json:"artist_id,omitempty"`
	ReleaseDate     *string `json:"release_date,omitempty"`
	AlbumArtPath    *string `json:"album_art_path,omitempty"`
	Genre           *string `json:"genre,omitempty"`
	Year            *int64  `json:"year,omitempty"`
	TotalTracks     *int64  `json:"total_tracks,omitempty"`
	CreatedAt       string  `json:"created_at"`
}

type Artist struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	NormalizedName string `json:"normalized_name"`
	CreatedAt      string `json:"created_at"`
}

type DownloadHistory struct {
	ID           int64   `json:"id"`
	UserID       *int64  `json:"user_id,omitempty"`
	TrackID      *int64  `json:"track_id,omitempty"`
	Quality      *int64  `json:"quality,omitempty"`
	Status       *string `json:"status,omitempty"`
	Service      *string `json:"service,omitempty"`
	StartedAt    string  `json:"started_at"`
	CompletedAt  *string `json:"completed_at,omitempty"`
	ErrorMessage *string `json:"error_message,omitempty"`
}

type Track struct {
	ID              int64   `json:"id"`
	Title           string  `json:"title"`
	NormalizedTitle string  `json:"normalized_title"`
	ArtistID        *int64  `json:"artist_id,omitempty"`
	AlbumID         *int64  `json:"album_id,omitempty"`
	Duration        *int64  `json:"duration,omitempty"`
	TrackNumber     *int64  `json:"track_number,omitempty"`
	DiscNumber      *int64  `json:"disc_number,omitempty"`
	SampleRate      *int64  `json:"sample_rate,omitempty"`
	BitDepth        *int64  `json:"bit_depth,omitempty"`
	Bitrate         *int64  `json:"bitrate,omitempty"`
	Channels        *int64  `json:"channels,omitempty"`
	Codec           *string `json:"codec,omitempty"`
	FilePath        string  `json:"file_path"`
	FileSize        *int64  `json:"file_size,omitempty"`
	ISRC            *string `json:"isrc,omitempty"`
	CreatedAt       string  `json:"created_at"`
}

type User struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Email     *string `json:"email,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	LastLogin *string `json:"last_login,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

type UserTrack struct {
	UserID       *int64 `json:"user_id,omitempty"`
	TrackID      *int64 `json:"track_id,omitempty"`
	SymlinkPath  string `json:"symlink_path"`
	DownloadDate string `json:"download_date"`
}

type StreamripSearchResult struct {
	Source    string `json:"source"`
	MediaType string `json:"media_type"`
	ID        string `json:"id"`
	Desc      string `json:"desc"`
	Data      struct {
		Title     string `json:"title"`
		Duration  int    `json:"duration"`
		Performer struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"performer"`
		Album struct {
			Title string `json:"title"`
			Image struct {
				Small     string `json:"small"`
				Thumbnail string `json:"thumbnail"`
				Large     string `json:"large"`
			} `json:"image"`
		} `json:"album"`
		ISRC string `json:"isrc"`
	} `json:"data"`
}

type TrackPreview struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Duration int    `json:"duration"`
	Image    string `json:"image"`
	TrackID  string `json:"track_id"`
	Source   string `json:"source"`
	ISRC     string `json:"isrc"`
}
